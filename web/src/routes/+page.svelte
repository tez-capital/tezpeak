<script lang="ts">
	import { nodes } from '@app/state/index';
	import {
		state as tezbakeStatus,
		bakers as tezbakeBakers,
		services as tezbakeServices,
		bakingRights,
		votingPeriodInfo
	} from '@app/state/tezbake';
	import NodeStatusCard from '@components/app/NodeStatusCard.svelte';
	import BakerStatusCard from '@components/app/BakerStatusCard.svelte';
	import BakerRightsCard from '@components/app/BakerRightsCard.svelte';
	import ServicesStatusCard from '@components/app/ServicesStatusCard.svelte';

	import GovernancePeriodCard from '@src/components/app/GovernancePeriodCard.svelte';
	import PayoutsCard from '@src/components/app/PayoutsCard.svelte';

	$: showBakerColors = $tezbakeBakers.length > 1;
	$: expandedBakingRights = $tezbakeBakers.length > 1;

	//$: services = $state.tezbake?.services;
	// $: hasAnyService =
	// 	Object.keys(services.node_services).length > 0 ||
	// 	Object.keys(services.signer_services).length > 0;

	// $: votingPeriodInfo = pickVotingPeriodInfo([bakerNode, ...nodes.map((n) => n[1])]);
	// $: votingPeriodBlock = getCurrentBlock([bakerNode, ...nodes.map((n) => n[1])]);
</script>

<div class="dashboard-grid-wrap">
	<div class="dashboard-grid">
		{#if Object.keys($tezbakeServices.applications ?? {}).length > 0}
			<ServicesStatusCard title="Baker's Services" services={$tezbakeServices} />
		{/if}

		{#if $tezbakeStatus}
			{#each $tezbakeBakers as [baker, info]}
				<BakerStatusCard baker={baker ?? {}} status={info} showColor={showBakerColors} />
			{/each}
			<GovernancePeriodCard votingPeriodInfo={$votingPeriodInfo} />
		{/if}
		<PayoutsCard votingPeriodInfo={$votingPeriodInfo} />
		{#each $nodes as [node, info]}
			<NodeStatusCard node={info} title={node} />
		{/each}
		{#if $tezbakeStatus}
			<div class="baker-rights" class:expanded={expandedBakingRights}>
				<BakerRightsCard
					mode="upcoming"
					rights={$bakingRights.future}
					{showBakerColors}
					title="Upcoming Baking Rights"
				/>
			</div>
			<div class="baker-rights" class:expanded={expandedBakingRights}>
				<BakerRightsCard
					mode="past"
					rights={$bakingRights.past}
					{showBakerColors}
					title="Past Baking Rights"
				/>
			</div>
		{/if}
	</div>
</div>

<style lang="sass">
.dashboard-grid-wrap
	display: grid
	grid-template-columns: 1fr minmax(0px, 1400px) 1fr
	width: calc(100% - var(--spacing) * 2)
	padding: var(--spacing)
	gap: var(--spacing)

	.dashboard-grid
		display: grid
		grid-column: 2
		grid-template-columns: minmax(450px, 1fr) minmax(450px, 1fr) minmax(450px, 1fr)
		gap: var(--spacing)

		.baker-rights
			display: grid
			grid-template-rows: 1fr

		.baker-rights.expanded
			grid-column: 1/-1

@media (max-width: 1400px)
	.dashboard-grid
		grid-template-columns: minmax(450px, 1fr) minmax(450px, 1fr) !important


@media (max-width: 900px) 
	.dashboard-grid
		grid-template-columns: minmax(450px, 1fr) !important
		

</style>
