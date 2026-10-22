<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button, Card } from '@alis-is/starlight';
	import type { VotingPeriodInfo } from '@src/common/types/status';
	import { getGovernancePeriodBlocksLeft, getVotingPeriodTimeLeft } from '@src/util/gov';
	import { now } from '@app/state/clock';

	interface Props {
		votingPeriodInfo: VotingPeriodInfo | undefined;
	}
	let { votingPeriodInfo }: Props = $props();
</script>

<Card>
	{#snippet header()}
		<h3 class="sl-card-title">Governance</h3>
		{#if votingPeriodInfo}<span class="hint">period #{votingPeriodInfo.voting_period.index}</span>{/if}
	{/snippet}

	{#if votingPeriodInfo}
		<div class="metric kind">{votingPeriodInfo.voting_period.kind}</div>
		<dl class="kv">
			<dt>Ends in</dt>
			<dd>{#key $now}{getVotingPeriodTimeLeft(votingPeriodInfo)}{/key}</dd>
			<dt>Blocks remaining</dt>
			<dd>{getGovernancePeriodBlocksLeft(votingPeriodInfo)}</dd>
		</dl>
	{:else}
		<div class="hint">No voting period data</div>
	{/if}

	{#snippet footer()}
		<Button onclick={() => goto('/governance')}>Open</Button>
	{/snippet}
</Card>

<style lang="sass">
.kind
	text-transform: capitalize
</style>
