<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button, Expander, ProgressBar } from '@alis-is/starlight';
	import FolderIcon from '@alis-is/starlight/la/icons/folder-solid.svelte';
	import HomeIcon from '@alis-is/starlight/la/icons/home-solid.svelte';
	import BackIcon from '@alis-is/starlight/la/icons/angle-left-solid.svelte';
	import { onMount } from 'svelte';
	import { listReports } from '@src/app/tezpay/client';

	let reports: Array<string> = [];
	let dryRunReports: Array<string> = [];
	let isLoading = true;

	$: sortedReports = reports.sort((a, b) => b.localeCompare(a));
	$: sortedDryRunReports = dryRunReports.sort((a, b) => b.localeCompare(a));

	onMount(async () => {
		try {
			isLoading = true;
			const promises = {
				reports: listReports(),
				dryRunReports: listReports(true)
			};
			await Promise.allSettled(Object.values(promises));
			reports = await promises.reports;
			dryRunReports = await promises.dryRunReports;
		} finally {
			isLoading = false;
		}
	});
</script>

<div class="reports-wrap">
	<nav class="page-navigation">
		<Button onclick={() => goto('/tezpay')}>
			<div class="navigation-btn-content"><BackIcon /> BACK</div>
		</Button>
		<Button onclick={() => goto('/')}>
			<div class="navigation-btn-content"><HomeIcon /> HOME</div>
		</Button>
	</nav>
	<div class="reports">
		{#if isLoading}
			<div class="progress-bar">
				<ProgressBar message="Loading..." progress="indeterminate" />
			</div>
		{:else}
			<Expander is_open={true}>
				{#snippet header_title()}
					<div class="title"><FolderIcon /> Reports</div>
				{/snippet}
				<div class="reports-list">
					{#if sortedReports.length === 0}
						<div class="no-data">NO DATA</div>
					{:else}
						{#each sortedReports as report}
							<Button label={report} onclick={() => goto(`/tezpay/reports/${report}`)} />
						{/each}
					{/if}
				</div>
			</Expander>
			<Expander>
				{#snippet header_title()}
					<div class="title"><FolderIcon /> Dry-Run Reports</div>
				{/snippet}
				<div class="reports-list">
					{#if sortedDryRunReports.length === 0}
						<div class="no-data">NO DATA</div>
					{:else}
						{#each sortedDryRunReports as report}
							<Button label={report} onclick={() => goto(`/tezpay/reports/${report}?dry=true`)} />
						{/each}
					{/if}
				</div>
			</Expander>
		{/if}
	</div>
</div>

<style lang="sass">
.reports-wrap
	display: grid
	gap: var(--spacing)

	.reports
		display: grid
		gap: var(--spacing)

	.title
		display: flex
		align-items: center
		gap: var(--spacing-f2)

		:global(svg)
			width: 24px
			height: 24px
			fill: var(--text-color)

.reports-list
	display: grid
	grid-template-columns: repeat(auto-fill, minmax(min(100%, 220px), 1fr))
	gap: var(--spacing-f2)
	padding: var(--spacing-f2) 0
	--button-vertical-spacing: var(--spacing-f2)
	--button-text-transform: none

.no-data
	grid-column: 1 / -1
	display: flex
	justify-content: center
	align-items: center
	font-size: 1.2rem
	color: var(--hint-color)
	margin: var(--spacing)

.progress-bar
	display: grid
	margin: var(--spacing-x2)
</style>
