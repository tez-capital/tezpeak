<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import {
		Button,
		Input,
		Card,
		Expander,
		ProgressBar,
		DataGrid,
		type DataGridColumn
	} from '@alis-is/starlight';
	import Separator from '@src/components/app/Separator.svelte';
	import BackIcon from '@alis-is/starlight/la/icons/angle-left-solid.svelte';

	import { onMount } from 'svelte';
	import { getReport } from '@src/app/tezpay/client';
	import type { PayoutReport, CycleReport } from '@src/common/types/tezpay';
	import { formatAddress, formatBalance, formatPercentage } from '@src/util/format';
	import { writeToClipboard } from '@src/util/clipboard';

	const reportId = $page.params.report ?? '';
	const dry = $page.url.searchParams.get('dry') === 'true';

	$: dryTitlePrefix = dry ? 'DRY ' : '';

	let loading = true;
	let error: string | undefined = undefined;
	let report: CycleReport | undefined = undefined;
	onMount(async () => {
		try {
			report = await getReport(reportId, dry);
		} catch (e: any) {
			error = `Error: ${e.message}`;
		} finally {
			loading = false;
		}
	});

	const COPYABLE: Array<keyof PayoutReport> = ['delegator', 'recipient', 'op_hash'];
	const columns: Array<DataGridColumn<PayoutReport>> = (
		[
			{ key: 'delegator', value: (item) => formatAddress(item.delegator) },
			{ key: 'delegator_balance', value: (item) => formatBalance(item.delegator_balance) },
			{ key: 'kind' },
			{ key: 'tx_kind' },
			{ key: 'recipient', value: (item) => formatAddress(item.recipient) },
			{ key: 'amount', value: (item) => formatBalance(item.amount) },
			{ key: 'fee_rate', value: (item) => formatPercentage(item.fee_rate), label: `Fee Rate` },
			{ key: 'fee', value: (item) => formatBalance(item.fee) },
			{ key: 'tx_fee', value: (item) => formatBalance(item.tx_fee), label: `Tx Fee` },
			{ key: 'op_hash', label: `Op Hash` },
			{ key: 'note' }
		] as Array<DataGridColumn<PayoutReport>>
	).map((c) => ({ sortable: true, width: 'minmax(max-content, 1fr)', ...c }));
</script>

{#snippet cell({ row, column }: { row: PayoutReport; column: DataGridColumn<PayoutReport> })}
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

<div class="report-wrap">
	<nav class="page-navigation">
		<Button onclick={() => goto('/tezpay/reports')}>
			<div class="navigation-btn-content"><BackIcon /> BACK</div>
		</Button>
	</nav>
	<div class="report">
		{#if loading}
			<ProgressBar message="Loading..." progress="indeterminate" />
		{:else if error}
			<div class="no-data error">{error}</div>
		{:else if !report}
			<div class="no-data">No data</div>
		{:else}
			<Card>
				<div class="report-grid">
					<div class="title">
						<h3>{dryTitlePrefix}Report #{report?.name}</h3>
					</div>
					<div class="summary">
						<div class="summary-section">
							<div class="summary-property">
								<Input
									label="Timestamp"
									value={report?.summary?.timestamp
										? new Date(report.summary.timestamp).toLocaleString()
										: '-'}
									readonly
								/>
							</div>
							<div class="summary-property">
								<Input label="delegators" value={report?.summary?.delegators ?? '-'} readonly />
							</div>
						</div>
						<div class="separator"><Separator /></div>
						<div class="title">
							<h5>Balances</h5>
						</div>
						<div class="summary-section">
							<div class="summary-property">
								<Input
									label="Own Staked"
									value={formatBalance(report?.summary?.own_staked_balance ?? '0')}
									readonly
								/>
							</div>
							<div class="summary-property">
								<Input
									label="Own Delegated"
									value={formatBalance(report?.summary?.own_delegated_balance ?? '0')}
									readonly
								/>
							</div>
							<div class="summary-property">
								<Input
									label="External Staked"
									value={formatBalance(report?.summary?.external_staked_balance ?? '0')}
									readonly
								/>
							</div>
							<div class="summary-property">
								<Input
									label="External Delegated"
									value={formatBalance(report?.summary?.external_delegated_balance ?? '0')}
									readonly
								/>
							</div>
						</div>
						<div class="separator"><Separator /></div>
						<div class="title">
							<h5>Fees & Income</h5>
						</div>
						<div class="summary-section">
							<div class="summary-property">
								<Input
									label="Cycle Fees"
									value={formatBalance(report?.summary?.cycle_fees ?? '0')}
									readonly
								/>
							</div>
							<div class="summary-property">
								<Input
									label="Cycle Rewards"
									value={formatBalance(report?.summary?.cycle_rewards ?? '0')}
									readonly
								/>
							</div>
							<div class="summary-property">
								<Input
									label="Distributed Rewards"
									value={formatBalance(report?.summary?.distributed_rewards ?? '0')}
									readonly
								/>
							</div>
							<div class="summary-property">
								<Input
									label="Bond Income"
									value={formatBalance(report?.summary?.bond_income ?? '0')}
									readonly
								/>
							</div>
							<div class="summary-property">
								<Input
									label="Fee Income"
									value={formatBalance(report?.summary?.fee_income ?? '0')}
									readonly
								/>
							</div>
							<div class="summary-property total">
								<Input
									label="Total Income"
									value={formatBalance(report?.summary?.total_income ?? '0')}
									readonly
								/>
							</div>
						</div>
						<div class="separator"><Separator /></div>
						<div class="title">
							<h5>Donations</h5>
						</div>
						<div class="summary-section">
							<div class="summary-property">
								<Input
									label="Donated Bonds"
									value={formatBalance(report?.summary?.donated_bonds ?? '0')}
									readonly
								/>
							</div>
							<div class="summary-property">
								<Input
									label="Donated Fees"
									value={formatBalance(report?.summary?.donated_fees ?? '0')}
									readonly
								/>
							</div>
							<div class="summary-property total">
								<Input
									label="Donated Total"
									value={formatBalance(report?.summary?.donated_total ?? '0')}
									readonly
								/>
							</div>
						</div>
					</div>
					<div class="separator"><Separator /></div>
					<Expander title="Valid Payouts ({report?.payouts?.length ?? 0})">
						<DataGrid data={report?.payouts ?? []} {columns} {cell} />
					</Expander>
					<Expander title="Invalid Payouts ({report?.invalid?.length ?? 0})">
						<DataGrid data={report?.invalid ?? []} {columns} {cell} />
					</Expander>
				</div>
			</Card>
		{/if}
	</div>
</div>

<style lang="sass">
.report-wrap
	display: grid
	gap: var(--spacing)

	.report
		display: grid
		min-width: 0

		.report-grid
			display: grid
			grid-template-columns: minmax(0, 1fr)
			gap: var(--spacing)
			grid-auto-rows: min-content

			.title
				display: flex
				justify-content: center
				text-align: center

				h3, h5
					margin: var(--spacing-f2) 0

			.summary
				display: grid
				gap: var(--spacing)
				--input-hint-font-size: 0rem
				--input-hint-vertical-spacing: 0rem

				.summary-section
					display: grid
					grid-template-columns: repeat(auto-fill, minmax(min(100%, 220px), 1fr))
					pointer-events: none
					gap: var(--spacing)

					.total
						--input-background-color: var(--accent-soft)
						grid-column: 1 / -1

.copyable
	cursor: pointer
	text-decoration: underline dotted
	text-underline-offset: 3px

.no-data
	display: flex
	justify-content: center
	align-items: center
	font-size: 1.2rem
	color: var(--hint-color)
	margin: var(--spacing)

.separator
	margin: var(--spacing-f2) 0

.error
	color: var(--error-color)
</style>
