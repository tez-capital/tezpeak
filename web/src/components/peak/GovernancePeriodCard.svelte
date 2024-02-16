<script lang="ts">
	import { goto } from '$app/navigation';
	import Separator from './Separator.svelte';
	import Card from '@components/starlight/components/Card.svelte';
	import type { VotingPeriodInfo } from '@src/common/types/status';
	import { getVotingPeriodTimeLeft } from '@src/util/gov';
	import { onDestroy } from 'svelte';

	export let votingPeriodInfo: VotingPeriodInfo | undefined;
	export let block: number | undefined;

	$: timeLeft = getVotingPeriodTimeLeft(votingPeriodInfo, block);

	const interval = setInterval(() => {
		timeLeft = getVotingPeriodTimeLeft(votingPeriodInfo, block);
	}, 500);

	onDestroy(() => clearInterval(interval));

	function open_governance() {
		goto('/governance');
	}
</script>

<button class="unstyle-button governance-wrap" on:click={() => open_governance()}>
	<Card class="governance-card">
		<div class="governance">
			<div class="title">
				<h5>Governance</h5>
			</div>
			<Separator />
			{#if votingPeriodInfo}
				<div class="period-info">
					<div class="kind">{votingPeriodInfo?.voting_period.kind}</div>
					<div class="period">period</div>
					<div class="index">
						#{votingPeriodInfo.voting_period.index}
					</div>
				</div>
				<Separator />

				<div class="remaining">
					ends in
					<div class="value">{timeLeft}</div>
				</div>
			{:else}
				<div class="no-data">NO DATA</div>
			{/if}
		</div>
	</Card>
</button>

<style lang="sass">
.governance-wrap
	display: grid
	grid-template-rows: 1fr
	width: 100%
	height: 100%
	user-select: none
	&:hover
		cursor: pointer
		transition: background-color 0.2s
		--card-background-color: #151515

	:global(.governance-card)
		box-sizing: border-box
		height: 100%

.governance
	display: grid
	grid-template-rows: auto auto 1fr auto auto
	height: 100%
	gap: var(--spacing)
	.title
		display: flex
		justify-content: center
		h5
			font-size: 1.5rem
			font-weight: 500
			margin: 0

	.period-info
		display: grid
		grid-template-rows: 1fr auto auto auto 1fr
		justify-content: center
		gap: var(--spacing-f2)
		
		.index, .kind
			display: flex
			justify-content: center
			font-size: 1.5rem
			font-weight: 500

		.index
			grid-row: 4

		.period 
			display: flex
			justify-content: center
			grid-row: 3

		.kind
			text-transform: uppercase
			grid-row: 2
	.remaining
		display: flex
		justify-content: right
		align-items: flex-end

		.value
			display: inline-block
			padding-left : var(--spacing-f2)
			font-size: 1.25rem
			font-weight: 500

	.no-data
		display: flex
		justify-content: center
		align-items: center
		font-size: 1.5rem
		font-weight: 500
		height: 100%
</style>
