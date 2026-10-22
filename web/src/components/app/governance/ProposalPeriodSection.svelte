<script lang="ts">
	import {
		AGORA_PROPOSAL_URL,
		EXPLORER_URL,
		PROPOSAL_PERIOD_QUORUM_PERCENTAGE
	} from '@src/common/constants';

	import { Button } from '@alis-is/starlight';
	import NewProposalDialog from './NewProposalDialog.svelte';
	import ThumbsUpIcon from '@alis-is/starlight/la/icons/thumbs-up-solid.svelte';
	import BookOpenIcon from '@alis-is/starlight/la/icons/book-open-solid.svelte';
	import SearchIcon from '@alis-is/starlight/la/icons/search-solid.svelte';

	import Separator from '../Separator.svelte';
	import { createEventDispatcher } from 'svelte';
	import { upvote_proposal } from '@src/app/governance';
	import type { TransactionEventDispatcher } from '@src/common/types/events';
	import type { ExtendedGovernanceProposalPeriodDetail } from '@src/common/types/governance';
	import { writeToClipboard } from '@src/util/clipboard';

	export let pkh: string | undefined;
	export let periodIndex: number;
	export let period: ExtendedGovernanceProposalPeriodDetail;

	let isNewProposalDialogOpen = false;

	const dispatch = createEventDispatcher<TransactionEventDispatcher>();

	function open_explorer(id: string) {
		window.open(`${EXPLORER_URL}/${id}`, '_blank');
	}

	function open_agora(id: string) {
		window.open(`${AGORA_PROPOSAL_URL}/${id}/${periodIndex}`, '_blank');
	}

	async function upvote(proposal: string) {
		if (!pkh) return;
		dispatch('tx_building', {
			message: `building proposal upvote - ${proposal}`,
			stage: 'building'
		});
		const result = await upvote_proposal(pkh, periodIndex, [proposal]);
		if (typeof result === 'string') {
			dispatch('tx_broadcasted', {
				message: `proposal upvote - ${proposal}`,
				opHash: result,
				stage: 'confirming'
			});
			return;
		}
		dispatch('error', {
			message: `Failed to upvote. Reason: ${result}`,
			error: result,
			stage: 'failed'
		});
	}
</script>

<div class="proposal-period-wrap">
	<div class="proposal-period">
		{#if period.proposals.length === 0}
			no proposals
		{:else}
			{#each period.proposals as proposal}
				<div class="proposal" class:success={proposal.quorumReached}>
					<div class="title">
						<button class="unstyle-button" on:click={() => writeToClipboard(proposal.proposal)}>
							{proposal.proposal}
						</button>
					</div>

					<div class="upvotes">
						{proposal.upvotes} upvotes ~ {proposal.percentage}%
						{#if proposal.quorumReached}
							({PROPOSAL_PERIOD_QUORUM_PERCENTAGE}% reached)
						{:else}
							(needs {PROPOSAL_PERIOD_QUORUM_PERCENTAGE}%)
						{/if}
					</div>

					<div class="upvote-bar">
						<svg width="100%" height="15">
							<rect width="100%" y="5" height="5" fill="var(--progress-bar-background-color)" rx="2" ry="2"></rect>
							<rect
								width={proposal.percentageString}
								y="5"
								height="5"
								fill="var(--progress-color)"
								rx="2"
								ry="2"
							></rect>
							<rect
								width="3"
								height="15"
								x={`${PROPOSAL_PERIOD_QUORUM_PERCENTAGE}%`}
								fill="var(--progress-color)"
								rx="2"
								ry="2"
							></rect>
						</svg>
					</div>

					<Button onclick={() => open_explorer(proposal.proposal)}>
						<div class="button-content">
							<div class="label">EXPLORER</div>
							<div class="icon"><SearchIcon /></div>
						</div>
					</Button>
					<Button onclick={() => open_agora(proposal.proposal)}>
						<div class="button-content">
							<div class="label">AGORA</div>
							<div class="icon"><BookOpenIcon /></div>
						</div>
					</Button>
					<div class="upvote-btn" class:disabled={!pkh}>
						{#if pkh && period.votes && period.votes[proposal.proposal].includes(pkh)}
							<div class="button-content upvoted">
								<div class="label">UPVOTED</div>
								<div class="icon"><ThumbsUpIcon /></div>
							</div>
						{:else}
							<Button onclick={() => upvote(proposal.proposal)}>
								<div class="button-content">
									<div class="label">UPVOTE</div>
									<div class="icon"><ThumbsUpIcon /></div>
								</div>
							</Button>
						{/if}
					</div>
				</div>
			{/each}
			<div class="separator">
				<Separator />
			</div>
		{/if}
	</div>
	<div class="new-proposal" class:disabled={!pkh}>
		<Button label="New Proposal" onclick={() => (isNewProposalDialogOpen = true)} />
	</div>
</div>

<NewProposalDialog
	{pkh}
	{periodIndex}
	bind:open={isNewProposalDialogOpen}
	on:tx_building={(event) => dispatch('tx_building', event.detail)}
	on:tx_broadcasted={(event) => dispatch('tx_broadcasted', event.detail)}
	on:error={(event) => dispatch('error', event.detail)}
/>

<style lang="sass">
	.proposal-period-wrap
		display: grid
		gap: var(--spacing)
		padding: 0
		--button-vertical-spacing: var(--spacing-f2)

		.proposal-period
			.proposal
				display: grid
				gap: var(--spacing-x2)
				padding: var(--spacing) 0
				grid-template-columns: minmax(100px, 1fr) minmax(100px, 1fr) minmax(100px, 1fr)
				--progress-color: var(--error-color)	
				

				&.success
					--progress-color: var(--success-color)		

				.title
					display: flex
					justify-content: center
					font-size: 1.25rem
					font-weight: 500
					grid-column: 1 / -1
					min-width: 0
					overflow-wrap: anywhere
					text-align: center
				
				.upvotes
					display: flex
					justify-content: center
					align-items: center
					font-size: 1.1rem
					grid-column: 1 / -1	

				.upvote-bar
					display: flex

					grid-column: 1 / -1	
					overflow: hidden	
				.upvote-btn
					--button-background-color: color-mix(in srgb, var(--success-color) 75%, #000)
					--button-hover-background-color: var(--success-color)
				
				.upvoted
					display: flex
					justify-content: center
					align-items: center
					width: 100%
					height: 100%
					color: var(--success-color)

			.separator
				padding-top: var(--spacing)
				padding-bottom: var(--spacing)
		//.new-proposal	
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

@media (max-width: 600px)
	.proposal-period-wrap .proposal-period .proposal
		grid-template-columns: minmax(0, 1fr)
		gap: var(--spacing)
</style>
