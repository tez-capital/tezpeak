<script lang="ts">
	import { Card, Overlay as OverlayDialog, Button, ProgressBar } from '@alis-is/starlight';
	import Separator from './Separator.svelte';
	import type { TransactionStage } from '@src/common/types/events';

	export let open: boolean;
	export let opHash: string | undefined;
	export let error: string | undefined;
	export let allowClose: boolean = true;
	export let stage: TransactionStage = 'building';

	function getProgressMessage() {
		switch (stage) {
			case 'building':
				return 'Building transaction...';
			case 'confirming':
				return 'Waiting for confirmation...';
			case 'applied':
				return 'CONFIRMED';
			case 'failed':
				return 'FAILED';
		}
	}

	$: progressMessage = getProgressMessage();
</script>

<OverlayDialog bind:open persistent={stage === 'building' || stage === 'confirming'}>
	<div class="center">
		<Card>
			<div class="transaction-confirmation">
				{#if error}
					<div class="error-message">{error}</div>
				{:else}
					<h5 class="title">
						Transaction {opHash}
					</h5>
					<Separator />
					{#if stage === 'applied'}
						<div class="status">CONFIRMED</div>
					{:else if stage === 'failed'}
						<div class="status failed">FAILED</div>
					{:else}
						<div class="progress">
							<ProgressBar message={progressMessage} progress="indeterminate" />
						</div>
						{#if stage === 'building'}
							<div class="note">
								NOTE: You may need to confirm the transaction with your signer.
							</div>
						{/if}
					{/if}
					<div class="explore" class:disabled={stage === 'building'}>
						<Button onclick={() => window.open(`https://better-call.dev/${opHash}`, '_blank')}>
							Open Explorer
						</Button>
					</div>
				{/if}
				{#if allowClose || stage === 'applied' || error}
					<Button label="close" onclick={() => (open = false)} />
				{/if}
			</div>
		</Card>
	</div>
</OverlayDialog>

<style lang="sass">
.center
	width: min(480px, 100%)

.transaction-confirmation
	display: grid
	gap: var(--spacing)
	align-items: center
	--button-vertical-spacing: var(--spacing-f2)

	.title
		display: flex
		justify-content: center
		font-size: 1.25rem
		font-weight: 500
		text-align: center
		overflow-wrap: anywhere
		margin: 0


	.status
		display: flex
		justify-content: center
		font-size: 1.5rem
		font-weight: 500
		color: var(--success-color)
		padding: var(--spacing-x2)

		&.failed
			color: var(--error-color)

	.progress 
		padding-top: var(--spacing-x2)

	.note
		color: var(--hint-color)
		text-align: center
		padding-bottom: var(--spacing-x2)

.error-message
	display: flex
	justify-content: center
	align-items: center
	color: var(--error-color)
	text-align: left
	white-space: pre-wrap
	padding: var(--spacing)
</style>
