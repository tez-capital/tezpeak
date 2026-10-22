<script lang="ts">
	import { goto } from '$app/navigation';
	import { Card, Button, Select, ProgressBar, type SelectItem } from '@alis-is/starlight';
	import ProposalPeriodSection from '@components/app/governance/ProposalPeriodSection.svelte';
	import ExpPromPeriodSection from '@components/app/governance/ExpPromPeriodSection.svelte';
	import TransactionDialog from '@components/app/TransactionDialog.svelte';
	import HomeIcon from '@alis-is/starlight/la/icons/home-solid.svelte';

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
	<nav class="page-navigation">
		<Button onclick={() => goto('/')}>
			<div class="navigation-btn-content"><HomeIcon /> HOME</div>
		</Button>
	</nav>
	<Card>
		{#snippet header()}
			<h3 class="sl-card-title">Voting key</h3>
			<span class="hint">{availablePkhsOptions.length} available</span>
		{/snippet}

		{#if availablePkhsOptions.length >= 1}
			<div class:disabled={availablePkhsOptions.length == 1}>
				<Select label="Vote as" bind:value={selectedPkh} options={availablePkhsOptions} />
			</div>
		{:else}
			<p class="hint">None of your keys can vote in this period.</p>
		{/if}

		{#snippet footer()}
			<span class="hint">Prefer a web wallet?</span>
			<Button label="Vote with TezGov" onclick={() => openTezgov()} />
		{/snippet}
	</Card>

	<Card>
		{#snippet header()}
			<h3 class="sl-card-title">
				{periodDetail ? `${periodDetail.info.voting_period.kind} period` : 'Voting period'}
			</h3>
			{#if periodDetail}<span class="hint">#{periodDetail.info.voting_period.index}</span>{/if}
		{/snippet}

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
				<dl class="kv">
					<dt>Ends</dt>
					<dd>{ended ? 'ended' : timeLeft}</dd>
					{#if !ended}
						<dt>Blocks remaining</dt>
						<dd>{remainingBlocks}</dd>
					{/if}
				</dl>
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
	grid-template-columns: minmax(0, 1fr)
	gap: var(--spacing)

.governance-period
	display: grid
	grid-template-columns: minmax(0, 1fr)
	gap: var(--spacing)

p
	margin: 0

.error-message
	display: flex
	justify-content: center
	align-items: center
	color: var(--error-color)
	text-align: center
	padding: var(--spacing)
</style>
