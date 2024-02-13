<script lang="ts">
	import { state } from '@app/state';
	import NodeStatusCard from '@components/peak/NodeStatusCard.svelte';
	import BakerStatusCard from '@components/peak/BakerStatusCard.svelte';
	import BakerRightsCard from '@components/peak/BakerRightsCard.svelte';
	import ServicesStatusCard from '@components/peak/ServicesStatusCard.svelte';
	import Separator from '@components/peak/Separator.svelte';
	import Card from '@components/starlight/components/Card.svelte';
	import { pickVotingPeriodInfo } from '@src/util/gov';
	import GovernancePeriodCard from '@src/components/peak/GovernancePeriodCard.svelte';

	$: bakerNode = $state.baker_node;

	$: nodes = Object.entries($state.nodes).sort((a, b) => (a[0] > b[0] ? 1 : -1));
	$: bakers = Object.entries($state.bakers.bakers).sort((a, b) => (a[0] > b[0] ? 1 : -1));
	$: showBakerColors = bakers.length > 1;
	$: expandedBakingRights = bakers.length > 1;

	$: upcomingRights = $state.rights.rights.filter((r) => r.level > $state.rights.level);
	$: pastRights = $state.rights.rights
		.filter((r) => r.level <= $state.rights.level)
		.sort((a, b) => (a.level < b.level ? 1 : -1));

	$: services = $state.services;
	$: hasAnyService =
		Object.keys(services.node_services).length > 0 ||
		Object.keys(services.signer_services).length > 0;

	$: votingPeriodInfo = pickVotingPeriodInfo([bakerNode, ...nodes.map((n) => n[1])]);
</script>

<div class="dashboard-grid-wrap">
	<div class="dashboard-grid">
		
		{#if hasAnyService}
			<ServicesStatusCard title="Baker's Services" {services} />
		{/if}
		<NodeStatusCard node={bakerNode} title="Baker's Node" />

		{#each bakers as [baker, info]}
			<BakerStatusCard {baker} status={info} showColor={showBakerColors} />
		{/each}

		<GovernancePeriodCard {votingPeriodInfo} />
		
		<div class="baker-rights" class:expanded={expandedBakingRights}>
			<BakerRightsCard
				mode="upcoming"
				rights={upcomingRights}
				{showBakerColors}
				title="Upcoming Baking Rights"
			/>
		</div>
		<div class="baker-rights" class:expanded={expandedBakingRights}>
			<BakerRightsCard
				mode="past"
				rights={pastRights}
				{showBakerColors}
				title="Past Baking Rights"
			/>
		</div>

		{#each nodes as [node, info]}
			<NodeStatusCard node={info} title={node} />
		{/each}
	</div>
</div>

<style lang="sass">
.dashboard-grid-wrap
	display: grid
	grid-template-columns: 1fr minmax(0px, 1400px) 1fr
	width: calc(100% - var(--spacing) * 2)
	padding: var(--spacing)

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
