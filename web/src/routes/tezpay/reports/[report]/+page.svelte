<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import Button from '@src/components/starlight/components/Button.svelte';
	import Input from '@src/components/starlight/components/Input.svelte';
	import Card from '@src/components/starlight/components/Card.svelte';
	import Expander from '@src/components/starlight/components/Expander.svelte';
	import DataTable from '@src/components/starlight/components/DataTable.svelte';
	import ProgressBar from '@src/components/starlight/components/ProgressBar.svelte';
	import Separator from '@src/components/app/Separator.svelte';
	import BackIcon from '@components/la/icons/angle-left-solid.svelte';

	import { onMount } from 'svelte';
	import { getReport } from '@src/app/tezpay/client';
	import type { PayoutReport, CycleReport } from '@src/common/types/tezpay';
	import { formatAddress, formatBalance, formatPercentage } from '@src/util/format';
	import type { DataTableColumn } from '@src/components/starlight/types';
	import { writeToClipboard } from '@src/util/clipboard';

	const reportId = $page.params.report;
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

	const columns: Array<DataTableColumn<PayoutReport>> = [
		{
			name: 'delegator',
			getValue: (item) => formatAddress(item.delegator),
			click: (item) => writeToClipboard(item.delegator)
		},
		{ name: 'delegator_balance', getValue: (item) => formatBalance(item.delegator_balance) },
		{ name: 'kind' },
		{ name: 'tx_kind' },
		{
			name: 'recipient',
			getValue: (item) => formatAddress(item.recipient),
			click: (item) => writeToClipboard(item.recipient)
		},
		{ name: 'amount', getValue: (item) => formatBalance(item.amount) },
		{ name: 'fee_rate', getValue: (item) => formatPercentage(item.fee_rate), label: `Fee Rate` },
		{ name: 'fee', getValue: (item) => formatBalance(item.fee) },
		{ name: 'tx_fee', getValue: (item) => formatBalance(item.fee), label: `Tx Fee` },
		{ name: 'op_hash', label: `Op Hash` },
		{ name: 'note' }
	];
</script>

<div class="report-wrap">
	<div class="navigation-wrap">
		<Button on:click={() => goto('/tezpay/reports')}>
			<div class="navigation-btn-content"><BackIcon /> BACK</div>
		</Button>
	</div>
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
					<Expander title="Valid Payouts">
						<DataTable data={report?.payouts ?? []} {columns} />
					</Expander>
					<Expander title="Invalid Payouts">
						<DataTable data={report?.invalid ?? []} {columns} />
					</Expander>
				</div>
			</Card>
		{/if}
	</div>
</div>

<style lang="sass">
.report-wrap
	display: grid
	grid-template-columns: 1fr minmax(0px, 1400px) 1fr
	width: calc(100% - var(--spacing) * 2)
	padding: var(--spacing)
	gap: var(--spacing)

	.report
		display: grid
		grid-column: 2

		.report-grid
			display: grid
			grid-template-columns: 1fr
			gap: var(--spacing)
			grid-auto-rows: min-content

			.title
				white-space: nowrap
				display: flex
				justify-content: center

			.summary
				display: grid
				gap: var(--spacing)
				--input-hint-font-size: 0rem
				--input-hint-vertical-spacing: 0rem

				.summary-section
					display: grid
					grid-template-columns: repeat(auto-fill, minmax(250px, 1fr))
					pointer-events: none

					gap: var(--spacing)
					
					*
						min-width: 250px

					.total
						--input-background-color: rgba(255,215,0, 0.12)
						grid-column: 1 / -1

	.title
		display: grid
		grid-template-columns: 25px auto
		gap: var(--spacing)
		fill: var(--text-color)
		align-items: center


.navigation-wrap 
	display: grid
	grid-template-columns: auto auto 1fr 
	grid-column: 2
	gap: var(--spacing)

	.navigation-btn-content
		display: grid
		grid-template-columns: auto auto
		align-items: center
		gap: var(--spacing)
		:global(svg)
			fill: var(--text-color)
			width: 30px
			height: 30px

.no-data
	display: flex
	justify-content: center
	align-items: center
	font-size: 1.2rem
	color: var(--text-color)
	margin: var(--spacing)

.separator
	margin: var(--spacing) 0
.error 
	color: var(--error-color)
</style>
