<script lang="ts">
	import { Overlay as OverlayDialog, Card, Button, ProgressBar, Expander, DataGrid, type DataGridColumn } from '@alis-is/starlight';
	import Terminal from '@alis-is/starlight/components/Terminal.svelte';
	import { generatePayuts, executePayuts } from '@src/app/tezpay/client';
	import { formatLogMessageForTerminal, type LogMessage } from '@src/util/log';
	import { PHASE_MAPPING } from '@src/common/phases';
	import type {
		BatchResults,
		ExecutionResult,
		PayoutBlueprint,
		PayoutRecipe
	} from '@src/common/types/tezpay';
	import { formatAddress, formatBalance, formatPercentage } from '@src/util/format';
	import { writeToClipboard } from '@src/util/clipboard';
	import TimesSolid from '@alis-is/starlight/la/icons/times-solid.svelte';
	import { status } from '@src/app/state/tezbake';

	let terminal: Terminal;
	let cycle: number | 'latest' = 0;
	let phase = 'Generating payouts...';
	let error: string | undefined = undefined;
	let stage: 'generate' | 'execute' = 'generate';
	let dry = false;

	let isOpen = false;
	let generatedPayouts: PayoutBlueprint | undefined = undefined;
	$: validPayouts = (generatedPayouts?.payouts ?? []).filter((x) => x.valid);
	$: invalidPayouts = (generatedPayouts?.payouts ?? []).filter((x) => !x.valid);
	// all valid payouts are selected by default, the user can deselect to skip them
	let selectedPayouts: Array<PayoutRecipe> = [];
	$: selectedPayouts = validPayouts;
	const payoutKey = (x: PayoutRecipe) => JSON.stringify(x);

	let batchResults: BatchResults | undefined = [];

	const COPYABLE = ['delegator', 'recipient', 'hash'];
	const toGridColumns = <T,>(columns: Array<DataGridColumn<T>>) =>
		columns.map((c) => ({ width: 'minmax(max-content, 1fr)', ...c }));

	const invalidColumns = toGridColumns<PayoutRecipe>([
		{ key: 'delegator', value: (item) => formatAddress(item.delegator) },
		{ key: 'cycle', value: (item) => `#${item.cycle}` },
		{ key: 'recipient', value: (item) => formatAddress(item.recipient) },
		{ key: 'delegator_balance', value: (item) => formatBalance(item.delegator_balance) },
		{ key: 'amount', value: (item) => formatBalance(item.amount) },
		{ key: 'fee_rate', value: (item) => formatPercentage(item.fee_rate) },
		{ key: 'fee', value: (item) => formatBalance(item.fee) },
		{ key: 'note' }
	]);

	const validColumns = toGridColumns<PayoutRecipe>([
		{ key: 'delegator', value: (item) => formatAddress(item.delegator) },
		{ key: 'cycle', value: (item) => `#${item.cycle}` },
		{ key: 'recipient', value: (item) => formatAddress(item.recipient) },
		{ key: 'kind' },
		{ key: 'delegator_balance', value: (item) => formatBalance(item.delegator_balance) },
		{ key: 'amount', value: (item) => formatBalance(item.amount) },
		{ key: 'fee_rate', value: (item) => formatPercentage(item.fee_rate) },
		{ key: 'fee', value: (item) => formatBalance(item.fee) }
	]);

	export async function generate(c: number | 'latest', dryRun: boolean = false) {
		isOpen = true;
		cycle = c;
		phase = 'Generating payouts...';
		stage = 'generate';
		dry = dryRun;
		terminal.clear();

		let executionResult: ExecutionResult | undefined = undefined;
		// undefined is latest
		try {
			await generatePayuts(
				cycle === 'latest' ? undefined : cycle,
				(msg) => {
					const event = JSON.parse(msg) as LogMessage;
					if (event.phase) {
						const k = event.phase as keyof typeof PHASE_MAPPING;
						phase = PHASE_MAPPING[k] ?? event.phase;
					}
					if (event.phase === 'result') {
						generatedPayouts = (event as { cycle_payout_blueprint?: PayoutBlueprint })
							.cycle_payout_blueprint;
					}
					if (event.phase !== 'execution_finished') {
						terminal.write(formatLogMessageForTerminal(event));
					} else {
						executionResult = event as unknown as ExecutionResult;
					}
				},
				dry
			);
			//@ts-ignore
			const exitCode = executionResult?.exit_code;
			if (exitCode !== 0) {
				throw new Error(`execution failed with exit code ${exitCode}`);
			}
		} catch (e: any) {
			error = `Failed to generate payouts: ${e.message}`;
		} finally {
			phase = 'execution_finished';
		}
	}

	type Batch = {
		id: string;
		hash: string;
		recipes: Array<PayoutRecipe>;
		status: string;
	};

	let batchStore: { [key: string]: boolean } = {};
	let batches: Array<Batch> = [];
	const batchColumns = toGridColumns<Batch>([
		{ key: 'id' },
		{ key: 'recipes', label: 'number of transactions', value: (item) => item.recipes.length },
		{ key: 'hash' },
		{ key: 'status' }
	]);

	async function execute_payouts() {
		if (!generatedPayouts) {
			return;
		}
		phase = 'Executing payouts...';
		stage = 'execute';
		const selected = new Set(selectedPayouts.map(payoutKey));
		const blueprint = {
			...generatedPayouts,
			payouts: [
				...validPayouts.filter((x) => selected.has(payoutKey(x))),
				...invalidPayouts,
				...validPayouts
					.filter((x) => !selected.has(payoutKey(x)))
					.map((x) => ({ ...x, valid: false, note: 'MANUALLY_SKIPPED' }))
			]
		};
		batches = [];

		let executionResult: ExecutionResult | undefined = undefined;
		try {
			await executePayuts(
				blueprint,
				(msg) => {
					const event = JSON.parse(msg) as LogMessage;
					if (event.phase) {
						const k = event.phase as keyof typeof PHASE_MAPPING;
						phase = PHASE_MAPPING[k] ?? event.phase;
						// TODO: improve, do not use any
						const data = event as any;
						switch (event.phase) {
							case 'executing_batch':
								if (batchStore[data.batch_id] === undefined) {
									batchStore[data.batch_id] = true;
									batches = [
										...batches,
										{
											id: data.batch_id,
											recipes: data.recipes,
											status: 'executing...',
											hash: ''
										}
									];
									///batches.push(batchStore[data.id]);
								}
							case 'batch_waiting_for_confirmation':
								batches = batches.map((x) => {
									if (x.id === data.batch_id) {
										return { ...x, status: 'waiting for confirmation', hash: data.op_hash };
									}
									return x;
								});
								break;
							case 'batch_execution_finished':
								if (data.error) {
									batches = batches.map((x) => {
										if (x.id === data.batch_id) {
											return { ...x, status: 'error' };
										}
										return x;
									});
								} else {
									batches = batches.map((x) => {
										if (x.id === data.batch_id) {
											return { ...x, status: 'success' };
										}
										return x;
									});
								}
								break;
						}
					}

					if (event.phase !== 'execution_finished') {
						terminal.write(formatLogMessageForTerminal(event));
					} else {
						executionResult = event as unknown as ExecutionResult;
					}
				},
				dry
			);
			//@ts-ignore
			const exitCode = executionResult?.exit_code;
			if (exitCode !== 0) {
				throw new Error(`execution failed with exit code ${exitCode}`);
			}
		} catch (e: any) {
			error = `Failed to execute payouts: ${e.message}`;
		} finally {
			phase = 'execution_finished';
		}
	}

	async function retry() {
		error = undefined;
		if (generatedPayouts === undefined) {
			await generate(cycle);
		} else {
			await execute_payouts();
		}
	}

	let isLogVisible = false;
	$: showLogButtonLabel = isLogVisible ? 'Hide log' : 'Show log';
	export function toggle_log() {
		isLogVisible = !isLogVisible;
	}
</script>

{#snippet cell({ row, column }: { row: any; column: DataGridColumn<any> })}
	{@const value = column.value ? column.value(row) : row[column.key]}
	{#if COPYABLE.includes(column.key) && row[column.key]}
		<button
			class="unstyle-button copyable"
			title="Copy to clipboard"
			onclick={() => writeToClipboard(String(row[column.key]))}>{value}</button
		>
	{:else}
		{value ?? '-'}
	{/if}
{/snippet}

<OverlayDialog bind:open={isOpen} persistent>
	<Card>
		{#snippet header()}
			<h3 class="sl-card-title">Payout · cycle {cycle}{dry ? ' · dry run' : ''}</h3>
			{#if phase === 'execution_finished'}
				<button class="unstyle-button close-icon" aria-label="Close" onclick={() => (isOpen = false)}>
					<TimesSolid />
				</button>
			{/if}
		{/snippet}
		<div class="payout-wrap">
			<div class="content">
				{#if stage === 'execute'}
					{#if phase !== 'execution_finished'}
						<div class="progress-wrap">
							<ProgressBar progress={'indeterminate'} message={phase} />
						</div>
					{:else if error}
						<div class="error">{error}</div>
						<Button onclick={retry} label="Retry"></Button>
					{:else}
						<div class="title">SUCCESS</div>
					{/if}
					{#if batches.length > 0}
						<div class="data"><DataGrid data={batches} columns={batchColumns} {cell} /></div>
					{:else if phase === 'execution_finished' && !error}
						<div class="title">Nothing to pay out.</div>
					{/if}
				{:else if stage === 'generate'}
					{#if phase !== 'execution_finished'}
						<div class="progress-wrap">
							<ProgressBar progress={'indeterminate'} message={phase} />
						</div>
					{:else if error}
						<div class="error">{error}</div>
						<Button onclick={retry} label="Retry"></Button>
					{:else}
						<Expander title="Invalid Payouts ({invalidPayouts.length})">
							<div class="data">
								<DataGrid data={invalidPayouts} columns={invalidColumns} {cell} />
							</div>
						</Expander>
						<Expander
							title="Valid Payouts ({selectedPayouts.length}/{validPayouts.length} selected)"
							is_open={true}
						>
							<div class="data">
								<DataGrid
									data={validPayouts}
									columns={validColumns}
									selectMode="multiple"
									rowKey={payoutKey}
									bind:selected={selectedPayouts}
									{cell}
								/>
							</div>
						</Expander>
					{/if}
				{:else}
					UNEXPECTED STATE
				{/if}
			</div>
			<div class="actions">
				<div class="secondary"><Button onclick={toggle_log} label={showLogButtonLabel}></Button></div>
				<span class="spacer"></span>
				{#if phase === 'execution_finished'}
					<div class="secondary"><Button onclick={() => (isOpen = false)} label="Close"></Button></div>
					{#if stage === 'generate' && !error && generatedPayouts}
						<Button onclick={execute_payouts} label="⚡ Execute"></Button>
					{/if}
				{/if}
			</div>
			<div class="terminal" class:hidden={!isLogVisible}>
				<Terminal bind:this={terminal} options={{ convertEol: true }}></Terminal>
			</div>
		</div>
	</Card>
</OverlayDialog>

<style lang="sass">
.payout-wrap
	position: relative
	box-sizing: border-box
	width: min(1200px, calc(100vw - var(--spacing) * 5))
	display: grid
	grid-template-columns: minmax(0, 1fr)
	gap: var(--spacing)
	--expanding-container-max-height: 45vh
	--expanding-container-overflow: auto
	--button-vertical-spacing: var(--spacing-f2)

	.title
		display: flex
		justify-content: center
		text-align: center

	.content
		display: grid
		grid-template-columns: minmax(0, 1fr)
		gap: var(--spacing)

		.data
			overflow: auto
			width: 100%

	.terminal
		overflow: auto
		height: 30vh

	.progress-wrap
		display: grid
		padding: var(--spacing-x2) 0

.actions
	display: flex
	flex-wrap: wrap
	align-items: center
	gap: var(--spacing-f2)
	--button-width: auto
	--button-horizontal-spacing: var(--spacing-x2)

	.spacer
		flex: 1

	.secondary
		--button-background-color: transparent
		--button-text-color: var(--text-color)
		--button-hover-background-color: var(--accent-soft)
		--button-hover-text-color: var(--text-color)
		--button-border: 1px solid var(--frost-border-color)

.copyable
	cursor: pointer
	text-decoration: underline dotted
	text-underline-offset: 3px

.error
	display: flex
	justify-content: center
	text-align: center
	color: var(--error-color)
	font-weight: 700
	font-size: 1.3rem
	padding: var(--spacing-x2)

.close-icon
	flex: none
	cursor: pointer
	fill: var(--text-color)
	width: 1.6rem
	height: 1.6rem
	transition: transform 0.2s

	&:hover
		transform: scale(1.1)
</style>
