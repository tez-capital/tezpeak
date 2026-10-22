<script lang="ts">
	import { nodes } from '@app/state/index';
	import {
		state as tezbakeStatus,
		bakers as tezbakeBakers,
		services as tezbakeServices,
		wallets as tezbakeWallets,
		votingPeriodInfo
	} from '@app/state/tezbake';
	import { state as tezpayStatus } from '@app/state/tezpay';
	import NodeStatusCard from '@components/app/NodeStatusCard.svelte';
	import BakerStatusCard from '@components/app/BakerStatusCard.svelte';
	import BakerRightsCard from '@components/app/BakerRightsCard.svelte';
	import ServicesStatusCard from '@components/app/ServicesStatusCard.svelte';
	import GovernancePeriodCard from '@components/app/GovernancePeriodCard.svelte';
	import PayoutsCard from '@components/app/PayoutsCard.svelte';
	import LedgerStatusCard from '@components/app/LedgerStatusCard.svelte';

	const showBakerColors = $derived($tezbakeBakers.length > 1);
	const hasServices = $derived(Object.keys($tezbakeServices.applications ?? {}).length > 0);
</script>

<div class="dashboard">
	{#if $tezbakeStatus && $tezbakeBakers.length}
		<h2 class="section-label">Bakers</h2>
		{#each $tezbakeBakers as [baker, info] (baker)}
			<BakerStatusCard {baker} status={info} showColor={showBakerColors} />
			<BakerRightsCard {baker} showColor={showBakerColors} />
		{/each}
	{/if}

	{#if $tezpayStatus || $tezbakeStatus}
		<h2 class="section-label">Operations</h2>
		{#if $tezpayStatus}<PayoutsCard />{/if}
		{#if $tezbakeStatus}
			<GovernancePeriodCard votingPeriodInfo={$votingPeriodInfo} />
			{#each $tezbakeWallets as [walletId, info] (walletId)}
				<LedgerStatusCard id={walletId} {info} />
			{/each}
		{/if}
	{/if}

	<h2 class="section-label">Infrastructure</h2>
	{#if hasServices}
		<ServicesStatusCard title="Baker services" services={$tezbakeServices} />
	{/if}
	{#each $nodes as [node, info] (node)}
		<NodeStatusCard node={info} title={node} />
	{/each}
</div>

<style lang="sass">
.dashboard
	display: grid
	grid-template-columns: repeat(auto-fill, minmax(min(100%, 320px), 1fr))
	gap: var(--spacing)
</style>
