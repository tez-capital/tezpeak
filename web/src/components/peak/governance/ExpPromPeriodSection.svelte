<script lang="ts">
	import { AGORA_PROPOSAL_URL, EXPLORER_URL } from '@src/common/constants';

	import Button from '@src/components/starlight/components/Button.svelte';
	import ThumbsUpIcon from '@src/components/la/icons/thumbs-up-solid.svelte';
	import ThumbsDownIcon from '@src/components/la/icons/thumbs-down-solid.svelte';
	import FistRaisedIcon from '@src/components/la/icons/fist-raised-solid.svelte';
	import BookOpenIcon from '@src/components/la/icons/book-open-solid.svelte';
	import SearchIcon from '@src/components/la/icons/search-solid.svelte';
	import Separator from '../Separator.svelte';

	import { createEventDispatcher } from 'svelte';
	import { cast_vote } from '@src/app/governance';
	import type { TransactionEventDispatcher } from '@src/common/types/events';
	import type {
		BallotVote,
		ExtendGovernanceExplorationOrPromotionPeriodDetail
	} from '@src/common/types/governance';
	import { writeToClipboard } from '@src/util/clipboard';

	export let pkh: string | undefined;
	export let periodIndex: number;
	export let period: ExtendGovernanceExplorationOrPromotionPeriodDetail;

	const dispatch = createEventDispatcher<TransactionEventDispatcher>();

	function open_explorer() {
		window.open(`${EXPLORER_URL}/${period.proposal}`, '_blank');
	}

	function open_agora() {
		window.open(`${AGORA_PROPOSAL_URL}/${period.proposal}/${periodIndex}`, '_blank');
	}

	async function vote(v: BallotVote) {
		if (!pkh) return;
		dispatch('tx_building', {
			message: `building proposal vote - ${period.proposal}`,
			stage: 'building'
		});
		const result = await cast_vote(pkh, periodIndex, period.proposal, v);
		if (typeof result === 'string') {
			dispatch('tx_broadcasted', {
				message: `proposal vote ${v} - ${period.proposal}`,
				opHash: result,
				stage: 'confirming'
			});
			return;
		}
		dispatch('error', {
			message: `Failed to vote. Reason: ${result}`,
			error: result,
			stage: 'failed'
		});
	}

	$: existingVote = period.ballots.find((b) => b.pkh === pkh);
</script>

<div
	class="period-wrap"
	class:quorum-reached={period.quorumReached}
	class:majority-reached={period.majorityReached}
>
	<div class="period">
		<div class="proposal">
			<div class="title">
				<button class="unstyle-button" on:click={() => writeToClipboard(period.proposal)}>
					{period.proposal}
				</button>
			</div>

			<div class="vote-summary">
				<div class="yay value">
					{period.summary.nay} ({period.summary.nayPercentage} %)
					<div class="icon"><ThumbsDownIcon /></div>
				</div>
				<div class="pass value">
					{period.summary.pass} ({period.summary.passPercentage} %)
					<div class="icon"><FistRaisedIcon /></div>
				</div>
				<div class="nay value">
					{period.summary.yayPercentage} ({period.summary.nayPercentage} %)
					<div class="icon"><ThumbsUpIcon /></div>
				</div>
			</div>
			<div class="separator">
				<Separator></Separator>
			</div>

			<div class="bar">
				<div class="bar-title">Quorum</div>
				<div class="bar-description">
					{period.quorumPercentage} %
					{#if period.quorumReached}
						({period.requiredQuorumPercentage}% reached)
					{:else}
						(needs {period.requiredQuorumPercentage}%)
					{/if}
				</div>
				<svg width="100%" height="15">
					<rect width="100%" height="5" y="5" fill="#e6e6e6" rx="2" ry="2"></rect>
					<rect
						width={`${period.summary.nayPercentage + period.summary.passPercentage + period.summary.yayPercentage}%`}
						y="5"
						height="5"
						fill="var(--yay-color)"
						rx="2"
						ry="2"
					></rect>
					<rect
						width={`${period.summary.nayPercentage + period.summary.passPercentage}%`}
						height="5"
						y="5"
						fill="var(--pass-color)"
						rx="2"
						ry="2"
					></rect>
					<rect
						width={`${period.summary.nayPercentage}%`}
						height="5"
						y="5"
						fill="var(--nay-color)"
						rx="2"
						ry="2"
					></rect>
					<rect
						width="3"
						height="15"
						x={`${period.requiredQuorumPercentage}%`}
						fill="var(--quorum-color)"
						rx="2"
						ry="2"
					></rect>
				</svg>
			</div>

			<div class="bar">
				<div class="bar-title">Supermajority</div>
				<div class="bar-description">
					{period.majorityPercentage}
					{#if period.majorityReached}
						({period.requiredMajorityPercentage}% reached)
					{:else}
						(required {period.requiredMajorityPercentage}%)
					{/if}
				</div>
				<svg width="100%" height="15">
					<rect width="100%" height="5" y="5" fill="#e6e6e6" rx="2" ry="2"></rect>
					<rect
						width={`${period.majorityPercentage}%`}
						y="5"
						height="5"
						fill="var(--majority-color)"
						rx="2"
						ry="2"
					></rect>
					<rect
						width="3"
						height="15"
						x={`${period.requiredMajorityPercentage}%`}
						fill="var(--majority-color)"
						rx="2"
						ry="2"
					></rect>
				</svg>
			</div>

			<Button on:click={() => open_explorer()}>
				<div class="button-content">
					<div class="label">EXPLORER</div>
					<div class="icon"><SearchIcon /></div>
				</div>
			</Button>
			<Button on:click={() => open_agora()}>
				<div class="button-content">
					<div class="label">AGORA</div>
					<div class="icon"><BookOpenIcon /></div>
				</div>
			</Button>
			<div class="separator">
				<Separator></Separator>
			</div>
			<div class="vote-section" class:disabled={!pkh}>
				{#if existingVote}
					<div class="existing-vote">You voted: <div class="ballot">{existingVote.ballot}</div></div>
				{:else}
					<div class="nay">
						<Button on:click={() => vote('nay')}>
							<div class="button-content">
								<div class="label">NAY</div>
								<div class="icon"><ThumbsDownIcon /></div>
							</div>
						</Button>
					</div>
					<div class="pass">
						<Button on:click={() => vote('pass')}>
							<div class="button-content">
								<div class="label">NAY</div>
								<div class="icon"><FistRaisedIcon /></div>
							</div>
						</Button>
					</div>
					<div class="yay">
						<Button on:click={() => vote('yay')}>
							<div class="button-content">
								<div class="label">NAY</div>
								<div class="icon"><ThumbsUpIcon /></div>
							</div>
						</Button>
					</div>
				{/if}
			</div>
		</div>
	</div>
</div>

<style lang="sass">
.period-wrap
	display: grid
	gap: var(--spacing)
	padding: var(--spacing)
	--yay-color: var(--success-color)
	--nay-color: var(--error-color)
	--pass-color: rgb(100, 100, 256)
	--quorum-color: var(--error-color)
	--majority-color: var(--error-color)

	&.quorum-reached
		--quorum-color: var(--success-color)

	&.majority-reached
		--pass-color: var(--success-color)
		--majority-color: var(--success-color)

	.period
		.proposal
			display: grid
			gap: var(--spacing)
			padding: var(--spacing)
			grid-template-columns: minmax(100px, 1fr) minmax(100px, 1fr) 	

			.title
				display: flex
				justify-content: center
				font-size: 1.25rem
				font-weight: 500
				grid-column: 1 / -1	

				div
					min-width: 0
					white-space: nowrap
					text-overflow: ellipsis
					overflow: hidden
			
			.vote-summary
				display: grid
				justify-content: center
				align-items: center
				grid-column: 1 / -1	
				gap: var(--spacing)
				padding: var(--spacing)
				grid-template-columns: auto auto auto

				.value
					display: flex
					justify-content: center
					align-items: center
					font-size: 0.8rem

					.icon
						padding-left: var(--spacing-f2)
						display: flex
						width: 20px
						fill: var(--button-text-color)


			.bar
				display: grid
				gap: var(--spacing)
				.bar-title
					display: flex
					justify-content: center
				.bar-description
					display: flex
					justify-content: center
					font-size: 0.8rem

			.vote-section
				display: grid
				grid-template-columns: 1fr 1fr 1fr
				align-items: center
				grid-column: 1 / -1	
				gap: var(--spacing)

				.existing-vote
					display: flex
					justify-content: center
					grid-column: 1 / -1	
					gap: var(--spacing-f2)

					.ballot
						font-weight: 500
						text-transform: uppercase

				.yay
					--button-background-color: rgb(0, 96, 0)
					--button-hover-background-color: rgb(0, 128, 0)
				.nay
					--button-background-color: rgb(96, 0, 0)
					--button-hover-background-color: rgb(128, 0, 0)
				.pass
					--button-background-color: rgb(0, 0, 96)
					--button-hover-background-color: rgb(0, 0, 128)

		.separator
			grid-column: 1 / -1
			padding-top: var(--spacing)
			padding-bottom: var(--spacing)

.button-content
	display: grid
	grid-template-columns: 1fr auto auto 1fr
	align-items: center
	gap: var(--spacing-f2)

	.label
		grid-column: 2

	.icon
		display: flex
		width: 20px
		fill: var(--button-text-color)
		align-items: center
</style>
