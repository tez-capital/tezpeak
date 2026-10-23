<script lang="ts">
	import { Card } from '@alis-is/starlight';
	import type { VotingPeriodInfo } from '@src/common/types/status';
	import { getVotingPeriodTimeLeft } from '@src/util/gov';
	import { now } from '@app/state/clock';

	interface Props {
		votingPeriodInfo: VotingPeriodInfo | undefined;
	}
	let { votingPeriodInfo }: Props = $props();

	const progress = $derived(
		votingPeriodInfo
			? votingPeriodInfo.position / (votingPeriodInfo.position + votingPeriodInfo.remaining || 1)
			: 0
	);
</script>

<Card>
	{#snippet header()}
		<a class="title-link" href="/governance"
			><h3 class="sl-card-title">Governance</h3>
			 →</a
		>
		{#if votingPeriodInfo}<span class="hint">period #{votingPeriodInfo.voting_period.index}</span
			>{/if}
	{/snippet}

	{#if votingPeriodInfo}
		<div class="metric kind">{votingPeriodInfo.voting_period.kind}</div>
		<div
			class="meter"
			role="progressbar"
			aria-valuenow={Math.round(progress * 100)}
			aria-label="voting period progress"
		>
			<span style:width="{progress * 100}%"></span>
		</div>
		<div class="hint">
			ends in {#key $now}{getVotingPeriodTimeLeft(votingPeriodInfo)}{/key} · {votingPeriodInfo.remaining.toLocaleString()}
			blocks
		</div>
	{:else}
		<div class="hint">No voting period data</div>
	{/if}
</Card>

<style lang="sass">
.kind
	text-transform: capitalize
</style>
