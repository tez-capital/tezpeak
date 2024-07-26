<script lang="ts">
	import { goto } from '$app/navigation';
	import Card from '@components/starlight/components/Card.svelte';
	import Button from '@components/starlight/components/Button.svelte';
	import Select from '@components/starlight/components/Select.svelte';
	import Separator from '@components/app/Separator.svelte';
	import ProposalPeriodSection from '@components/app/governance/ProposalPeriodSection.svelte';
	import ExpPromPeriodSection from '@components/app/governance/ExpPromPeriodSection.svelte';
	import TransactionDialog from '@components/app/TransactionDialog.svelte';
	import HomeIcon from '@components/la/icons/home-solid.svelte';

	import {
		getGovernancePeriodBlocksLeft,
		getVotingPeriodEndDate,
		preprocessExpPromPeriodDetail,
		preprocessProposalPeriodDetailProposals
	} from '@src/util/gov';
	import { formatDistanceToNow } from 'date-fns';
	import { onDestroy, onMount } from 'svelte';
	import axios from 'axios';
	import type {
		GovernancePeriodDetail,
		GovernanceProposalPeriodDetail,
		GovernanceExplorationOrPromotionPeriodDetail,
		CommonPeriodDetail
	} from '@src/common/types/governance';
	import type { SelectItem } from '@src/components/starlight/types';
	import ProgressBar from '@src/components/starlight/components/ProgressBar.svelte';
	import CommonPeriodSection from '@src/components/app/governance/CommonPeriodSection.svelte';
	import type {
		TransactionErrorEvent,
		TransactionBroadcastedEvent,
		TransactionStage,
		TransactionBuildingEvent
	} from '@src/common/types/events';
	import { openTezgov, waitConfirmation } from '@src/app/governance';

	let canVote: boolean | undefined;
	let periodDetail: GovernancePeriodDetail;
	let availablePkhs: Array<string> = [];
	let selectedPkh: SelectItem | undefined;

	$: voters = periodDetail?.voters ?? [];
	$: availablePkhsOptions = (availablePkhs ?? [])
		.filter((pkh) => voters.find((voter) => voter.pkh === pkh))
		.map((pkh) => ({ value: pkh, label: pkh }));

	$: proposalPeriod =
		periodDetail?.info.voting_period.kind === 'proposal'
			? preprocessProposalPeriodDetailProposals(periodDetail as GovernanceProposalPeriodDetail)
			: undefined;
	$: explorationOrPromotionPeriod = ['exploration', 'promotion'].includes(
		periodDetail?.info.voting_period.kind
	)
		? preprocessExpPromPeriodDetail(periodDetail as GovernanceExplorationOrPromotionPeriodDetail)
		: undefined;

	$: commonPeriod = ['testing', 'adoptions'].includes(periodDetail?.info.voting_period.kind)
		? (periodDetail as CommonPeriodDetail)
		: undefined;

	$: endDate = getVotingPeriodEndDate(periodDetail?.info);
	$: timeLeft = formatDistanceToNow(endDate, { addSuffix: true });

	$: ended = endDate < new Date();
	$: remainingBlocks = getGovernancePeriodBlocksLeft(periodDetail?.info);

	const interval = setInterval(() => {
		remainingBlocks = getGovernancePeriodBlocksLeft(periodDetail?.info);
		timeLeft = formatDistanceToNow(endDate, { addSuffix: true });
		ended = endDate < new Date();
	}, 500);

	onDestroy(() => clearInterval(interval));

	let loading: boolean = true;
	let error: undefined | string;
	async function fetchCanVote() {
		const result = await axios.get('/api/governance/can-vote');
		canVote = result.data as boolean;
	}
	async function fetchGovernancePeriodDetail() {
		if (!canVote) return;
		const result = await axios.get('/api/governance/period-detail');
		periodDetail = result.data;
	}
	async function fetchAvailablePkhs() {
		const result = await axios.get('/api/governance/available-pkhs');
		availablePkhs = result.data ?? [];
		if (!selectedPkh || availablePkhs.indexOf(selectedPkh.value) === -1)
			selectedPkh = { value: availablePkhs[0], label: availablePkhs[0] };
	}

	async function refreshData() {
		if (!canVote) return;
		try {
			loading = true;
			const promises = {
				detail: fetchGovernancePeriodDetail(),
				availablePkhs: fetchAvailablePkhs()
			};
			await Promise.allSettled(Object.values(promises));
			await promises.detail;

			// availablePkhs = [
			// 	'tz1P6WKJu2rcbxKiKRZHKQKmKrpC9TfW1AwM'
			// 	//'tz1hZvgjekGo7DmQjWh7XnY5eLQD8wNYPczE'
			// ];
			// selectedPkh = { value: availablePkhs[0], label: availablePkhs[0] };
		} catch (err: any) {
			if (err instanceof Error) {
				error = err.message;
			} else {
				error = err;
			}
		} finally {
			loading = false;
		}
	}

	onMount(async () => {
		try {
			await fetchCanVote();
			await refreshData();
		} finally {
			loading = false;
		}
	});

	let actionError: string | undefined;
	let transactionStage: TransactionStage = 'building';
	let opHash: string | undefined;
	let isTransactionDialogOpen: boolean = false;

	$: if (!isTransactionDialogOpen) {
		actionError = undefined;
		opHash = undefined;
		transactionStage = 'building';
	}

	function handle_tx_building(event: CustomEvent<TransactionBuildingEvent>) {
		actionError = undefined;
		opHash = undefined;
		isTransactionDialogOpen = true;
		transactionStage = event.detail.stage;
	}

	function handle_action_error(event: CustomEvent<TransactionErrorEvent>) {
		actionError = event.detail.message;
		isTransactionDialogOpen = true;
		transactionStage = event.detail.stage;
	}
	async function handle_tx_broadcasted(event: CustomEvent<TransactionBroadcastedEvent>) {
		opHash = event.detail.opHash;
		isTransactionDialogOpen = true;
		transactionStage = event.detail.stage;
		const result = await waitConfirmation(opHash);
		if (typeof result === 'boolean') {
			transactionStage = result ? 'applied' : 'failed';
		} else {
			actionError = result.message;
		}
	}
</script>

<div class="governance-wrap">
	<div class="navigation-wrap">
		<Button on:click={() => goto('/')}>
			<div class="navigation-btn-content"><HomeIcon /> HOME</div>
		</Button>
	</div>
	<div class="vote-as">
		<Card>
			<div class="vote-as-content">
				{#if availablePkhsOptions.length >= 1}
					<div class="vote-as-select" class:disabled={availablePkhsOptions.length == 1}>
						<Select label="Vote As" bind:value={selectedPkh} options={availablePkhsOptions} />
					</div>
					<div class="open-tezgov">
						<Button label="or vote with TEZGOV" on:click={() => openTezgov()} />
					</div>
				{:else if availablePkhsOptions.length === 0}
					<div class="no-pkhs">❗You have no available PKH to vote with.❗</div>
					<div class="open-tezgov">
						<Button label="vote with TEZGOV" on:click={() => openTezgov()} />
					</div>
				{/if}
			</div>
		</Card>
	</div>

	<div class="governance">
		<Card>
			<div class="governance-period">
				{#if loading}
					<ProgressBar progress="indeterminate" />
				{:else if error}
					<div class="error-message">
						{error}
					</div>
				{:else if !canVote}
					<div class="error-message">You are not eligible to vote</div>
				{:else}
					<div class="title">
						<h5>
							{periodDetail?.info.voting_period.kind} Period
						</h5>
					</div>
					<div class="governance-period-header">
						<div class="governance-period-number">
							#{periodDetail.info.voting_period.index}
						</div>
						<div class="governance-period-remaining">
							{#if ended}
								ended
							{:else}
								ends in {remainingBlocks} blocks ~ {timeLeft}
							{/if}
						</div>
					</div>
					<Separator />
					{#if proposalPeriod !== undefined}
						<ProposalPeriodSection
							period={proposalPeriod}
							periodIndex={periodDetail.info.voting_period.index}
							pkh={selectedPkh?.value}
							on:tx_building={handle_tx_building}
							on:error={handle_action_error}
							on:tx_broadcasted={handle_tx_broadcasted}
						/>
					{:else if explorationOrPromotionPeriod !== undefined}
						<ExpPromPeriodSection
							period={explorationOrPromotionPeriod}
							periodIndex={periodDetail.info.voting_period.index}
							pkh={selectedPkh?.value}
							on:tx_building={handle_tx_building}
							on:error={handle_action_error}
							on:tx_broadcasted={handle_tx_broadcasted}
						/>
					{:else if commonPeriod !== undefined}
						<CommonPeriodSection
							period={commonPeriod}
							periodIndex={periodDetail.info.voting_period.index}
						/>
					{:else}
						<div class="error-message">
							{periodDetail.info.voting_period.kind} is not supported yet
						</div>
					{/if}
				{/if}
			</div>
		</Card>
	</div>
</div>

<TransactionDialog
	bind:open={isTransactionDialogOpen}
	{opHash}
	error={actionError}
	allowClose={false}
	stage={transactionStage}
/>

<style lang="sass">
.governance-wrap
	display: grid
	grid-template-columns: 1fr minmax(0px, 1400px) 1fr
	width: calc(100% - var(--spacing) * 2)
	padding: var(--spacing)
	--button-vertical-spacing: var(--spacing-f2)
	gap: var(--spacing)

	.vote-as
		padding: var(--spacing)
		grid-column: 2
		
		.vote-as-content
			display: grid
			grid-template-columns: 1fr 1fr
			--button-text-transform: none

			.no-pkhs
				display: flex
				justify-content: center
				padding: var(--spacing)
				grid-column: 1 / -1
				
			.vote-as-select
				display: flex
				justify-content: center
				align-items: center
				height: 100%
				color: var(--error-color)
				grid-column: 1 / -1

			.open-tezgov
				display: flex
				justify-content: end
				align-items: center
				padding-top: var(--spacing)
				color: var(--error-color)
				grid-column: 2

	.governance
		display: grid
		grid-column: 2
		grid-template-columns: 1fr
		gap: var(--spacing)
		padding: var(--spacing)

		.governance-period
			display: grid
			gap: var(--spacing)
			padding: var(--spacing)
			grid-template-columns: minmax(200px, 1fr)

			.title
				display: flex
				justify-content: center
				position: relative
				text-transform: capitalize

				h5
					font-size: 1.5rem
					font-weight: 500
					margin: 0

			.governance-period-header
				display: grid
				grid-template-columns: auto 1fr auto
				.governance-period-remaining
					grid-column: 3

				.governance-period-number
					grid-column: 1

.navigation-wrap 
	display: grid
	grid-template-columns: auto 1fr
	grid-column: 2

	.navigation-btn-content
		display: grid
		grid-template-columns: auto auto
		align-items: center
		gap: var(--spacing)
		:global(svg)
			fill: var(--text-color)
			width: 30px
			height: 30px

.error-message
	display: flex
	justify-content: center
	align-items: center
	color: var(--error-color)
	text-align: center
	padding: var(--spacing)


</style>
