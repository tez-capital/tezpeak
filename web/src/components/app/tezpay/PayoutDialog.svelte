<script lang="ts">
	import OverlayDialog from '@components/starlight/components/OverlayDialog.svelte';
	import Card from '@components/starlight/components/Card.svelte';
	import Terminal from '@components/starlight/components/Terminal.svelte';
	import Button from '@components/starlight/components/Button.svelte';
	import ProgressBar from '@components/starlight/components/ProgressBar.svelte';
	import { DataTableSelectMode, type DataTableColumn } from '@components/starlight/types';
	import { generatePayuts, executePayuts } from '@src/app/tezpay/client';
	import { formatLogMessageForTerminal, type LogMessage } from '@src/util/log';
	import { PHASE_MAPPING } from '@src/common/phases';
	import type {
		BatchResults,
		ExecutionResult,
		PayoutBlueprint,
		PayoutRecipe
	} from '@src/common/types/tezpay';
	import DataTable from '@src/components/starlight/components/DataTable.svelte';
	import Expander from '@src/components/starlight/components/Expander.svelte';
	import { formatAddress, formatBalance, formatPercentage } from '@src/util/format';
	import { writeToClipboard } from '@src/util/clipboard';
	import TimesSolid from '@src/components/la/icons/times-solid.svelte';
	import { status } from '@src/app/state/tezbake';

	let terminal: Terminal;
	let cycle: number | 'latest' = 0;
	let phase = 'Generating payouts...';
	let error: string | undefined = undefined;
	let stage: 'generate' | 'execute' = 'generate';

	let isOpen = false;
	let generatedPayouts: PayoutBlueprint | undefined = undefined;
	$: generatedPayoutsItems = (generatedPayouts?.payouts ?? []).map((x) => ({
		...x,
		selected: x.valid
	}));

	$: validPayouts = generatedPayoutsItems.filter((x) => x.valid);
	$: invalidPayouts = generatedPayoutsItems.filter((x) => !x.valid);

	let batchResults: BatchResults | undefined = [];

	const invalidColumns: Array<DataTableColumn<PayoutRecipe>> = [
		{
			name: 'delegator',
			getValue: (item) => formatAddress(item.delegator),
			click: (item) => writeToClipboard(item.delegator)
		},
		{ name: 'cycle', getValue: (item) => `#${item.cycle}` },
		{
			name: 'recipient',
			getValue: (item) => formatAddress(item.recipient),
			click: (item) => writeToClipboard(item.recipient)
		},
		{ name: 'delegator_balance', getValue: (item) => formatBalance(item.delegator_balance) },
		{ name: 'amount', getValue: (item) => formatBalance(item.amount) },
		{ name: 'fee_rate', getValue: (item) => formatPercentage(item.fee_rate) },
		{ name: 'fee', getValue: (item) => formatBalance(item.fee) },
		{ name: 'note' }
	];

	const validColumns: Array<DataTableColumn<PayoutRecipe>> = [
		{
			name: 'delegator',
			getValue: (item) => formatAddress(item.delegator),
			click: (item) => writeToClipboard(item.delegator)
		},
		{ name: 'cycle', getValue: (item) => `#${item.cycle}` },
		{
			name: 'recipient',
			getValue: (item) => formatAddress(item.recipient),
			click: (item) => writeToClipboard(item.recipient)
		},
		{ name: 'kind' },
		{ name: 'delegator_balance', getValue: (item) => formatBalance(item.delegator_balance) },
		{ name: 'amount', getValue: (item) => formatBalance(item.amount) },
		{ name: 'fee_rate', getValue: (item) => formatPercentage(item.fee_rate) },
		{ name: 'fee', getValue: (item) => formatBalance(item.fee) }
	];

	export async function generate(c: number | 'latest') {
		isOpen = true;
		cycle = c;
		phase = 'Generating payouts...';
		stage = 'generate';
		let executionResult: ExecutionResult | undefined = undefined;
		// undefined is latest
		try {
			await generatePayuts(cycle === 'latest' ? undefined : cycle, (msg) => {
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
			});
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
	const batchColumns: Array<DataTableColumn<Batch>> = [
		{ name: 'id' },
		{
			name: 'recipes',
			label: 'number of transactions',
			getValue: (item) => item.recipes.length
		},
		{ name: 'hash', click: (item) => writeToClipboard(item.hash) },
		{ name: 'status' }
	];

	async function execute_payouts() {
		if (!generatedPayouts) {
			return;
		}
		phase = 'Executing payouts...';
		stage = 'execute';
		const blueprint = {
			...generatedPayouts,
			payouts: [
				...validPayouts.filter((x) => x.selected),
				...invalidPayouts,
				...validPayouts
					.filter((x) => !x.selected)
					.map((x) => ({ ...x, valid: false, note: 'MANUALLY_SKIPPED' }))
			]
		};
		batches = [];

		let executionResult: ExecutionResult | undefined = undefined;
		try {
			await executePayuts(blueprint, (msg) => {
				const event = JSON.parse(msg) as LogMessage;
				if (event.phase) {
					const k = event.phase as keyof typeof PHASE_MAPPING;
					phase = PHASE_MAPPING[k] ?? event.phase;
					// TODO: improve, do not use any
					const data = event as any;
					console.log(event.phase, event);
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
			});
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

<OverlayDialog bind:open={isOpen} persistent>
	<Card>
		<div class="payout-wrap">
			{#if phase === 'execution_finished'}
				<button class="unstyle-button close-icon" on:click={() => (isOpen = false)}>
					<TimesSolid />
				</button>
			{/if}
			<div class="title">
				<h3>Payout for cycle #{cycle}</h3>
			</div>
			<div class="content">
				{#if stage === 'execute'}
					{#if phase !== 'execution_finished'}
						<div class="progress-wrap">
							<div />
							<ProgressBar progress={'indeterminate'} message={phase} />
						</div>
					{:else if error}
						<div class="error">{error}</div>
						<Button on:click={retry} label="Retry"></Button>
					{:else}
						<div class="title">SUCCESS</div>
					{/if}
					{#if batches.length > 0}
						<DataTable data={batches} columns={batchColumns}></DataTable>
					{/if}
				{:else if stage === 'generate'}
					{#if phase !== 'execution_finished'}
						<div class="progress-wrap">
							<div />
							<ProgressBar progress={'indeterminate'} message={phase} />
						</div>
					{:else if error}
						<div class="error">{error}</div>
						<Button on:click={retry} label="Retry"></Button>
					{:else}
						<Expander title="Invalid Payouts">
							<div class="data">
								<DataTable
									data={invalidPayouts}
									columns={invalidColumns}
									selectMode={DataTableSelectMode.none}
								></DataTable>
							</div>
						</Expander>
						<Expander title="Valid Payouts" isOpen={true}>
							<div class="data">
								<DataTable
									data={validPayouts}
									columns={validColumns}
									selectMode={DataTableSelectMode.multiple}
								></DataTable>
							</div>
						</Expander>
						<Button on:click={execute_payouts} label="⚡ Execute ⚡"></Button>
					{/if}
				{:else}
					UNEXPECTED STATE
				{/if}
			</div>
			{#if phase === 'execution_finished'}
				<Button on:click={() => (isOpen = false)} label="close"></Button>
			{/if}
			<Button on:click={toggle_log} label={showLogButtonLabel}></Button>
			<div class="terminal" class:hidden={!isLogVisible}>
				<Terminal bind:this={terminal} options={{ convertEol: true }}></Terminal>
			</div>
		</div>
	</Card>
</OverlayDialog>

<style lang="sass">
.payout-wrap
	position: relative
	width: 90vw
	height: auto
	display: grid
	grid-template-: auto 1fr auto
	gap: var(--spacing-x2)
	--expanding-container-max-height: 25vh
	--expanding-container-overflow: auto
	
	.title
		display: flex
		justify-content: center
		align-items: center

	.content
		max-height: 60vh
		max-width: 100%
		display: grid
		grid-auto-rows: auto
		gap: var(--spacing)

		.data
			overflow: auto
			width: 100%
			height: auto

	.terminal
		overflow: auto
		height: 30vh
	
	.progress-wrap
		display: grid
		height: 100%
		grid-template-rows: 1fr auto 1fr 

.error
	display: flex
	justify-content: center
	color: var(--error-color)
	font-weight: 700
	font-size: 1.3rem
	padding: var(--spacing-x2)

.close-icon
	position: absolute
	top: var(--spacing)
	right: var(--spacing)
	background: transparent
	border: none
	cursor: pointer
	fill: var(--text-color)
	width: 2rem
	height: 2rem
	transitions: 0.2s

	&:hover
		transform: scale(1.1)
</style>
