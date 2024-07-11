<script lang="ts">
	import { goto } from '$app/navigation';
	import Button from '@components/starlight/components/Button.svelte';
	import FolderIcon from '@components/la/icons/folder-solid.svelte';
	import Expander from '@components/starlight/components/Expander.svelte';
	import HomeIcon from '@components/la/icons/home-solid.svelte';
	import BackIcon from '@components/la/icons/angle-left-solid.svelte';
	import ProgressBar from '@src/components/starlight/components/ProgressBar.svelte';
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
	<div class="navigation-wrap">
		<Button on:click={() => goto('/tezpay')}>
			<div class="navigation-btn-content"><BackIcon /> BACK</div>
		</Button>
		<Button on:click={() => goto('/')}>
			<div class="navigation-btn-content"><HomeIcon /> HOME</div>
		</Button>
	</div>
	<div class="reports">
		{#if isLoading}
			<div class="progress-bar">
				<ProgressBar message="Loading..." progress="indeterminate" />
			</div>
		{:else}
			<Expander isOpen={true}>
				<div class="title" slot="header-title">
					<FolderIcon /> Reports
				</div>
				<div class="reports-list">
					{#if sortedReports.length === 0}
						<div class="no-data">NO DATA</div>
					{:else}
						{#each sortedReports as report}
							<Button label={report} on:click={() => goto(`/tezpay/reports/${report}?dry=true`)} />
						{/each}
					{/if}
				</div>
			</Expander>
			<Expander>
				<div class="title" slot="header-title">
					<FolderIcon /> Dry-Run Reports
				</div>
				<div class="reports-list">
					{#if sortedDryRunReports.length === 0}
						<div class="no-data">NO DATA</div>
					{:else}
						{#each sortedDryRunReports as report}
							<Button label={report} on:click={() => goto(`/tezpay/reports/${report}?dry=true`)} />
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
	grid-template-columns: 1fr minmax(0px, 1400px) 1fr
	width: calc(100% - var(--spacing) * 2)
	padding: var(--spacing)
	gap: var(--spacing)

	.reports
		display: grid
		grid-column: 2

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
			wdith: 30px
			height: 30px

.reports-list
	display: grid
	grid-template-columns: 1fr
	gap: var(--spacing)

.no-data
	display: flex
	justify-content: center
	align-items: center
	font-size: 1.2rem
	color: var(--text-color)
	margin: var(--spacing)

.progress-bar
	display: grid
	margin: var(--spacing-x2)
</style>
