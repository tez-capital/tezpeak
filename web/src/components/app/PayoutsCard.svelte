<script lang="ts">
	import { goto } from '$app/navigation';
	import Card from '@components/starlight/components/Card.svelte';
	import Button from '../starlight/components/Button.svelte';
	import Separator from './Separator.svelte';
	import { services, wallet } from '@app/state/tezpay';
	import { formatBalance } from '@src/util/format';

	function open_governance() {
		goto('/tezpay');
	}

	$: isTezpayRunning = $services.applications?.tezpay?.['tezpay']?.status === 'running';
</script>

<div class="governance-wrap">
	<Card class="governance-card">
		<div class="governance">
			<div class="title">
				<h5>PAYOUTS</h5>
			</div>
			<Separator />

			<div class="payouts-info">
				<div class="row" />
				<div class="property">Automatic Payouts:</div>
				{#if isTezpayRunning}
					<div class="value automatic-payouts-status ok">ACTIVE</div>
				{:else}
					<div class="value automatic-payouts-status warn">INACTIVE</div>
				{/if}

				{#if $wallet !== undefined}
					<div class="property">Wallet Balance:</div>
					<div
						class="value balance"
						class:error={$wallet.level === 'error'}
						class:warn={$wallet.level === 'warning'}
					>
						{formatBalance($wallet.balance)}
					</div>
				{/if}
				<div class="row" />
			</div>
			<Separator />
			<div class="enter">
				<Button on:click={open_governance}>OPEN</Button>
			</div>
		</div>
	</Card>
</div>

<style lang="sass">
.governance-wrap
	display: grid
	grid-template-rows: 1fr
	width: 100%
	height: 100%
	user-select: none
	// &:hover
	// 	cursor: pointer
	// 	transition: background-color 0.2s
	// 	--card-background-color: #151515

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

	.payouts-info
		display: grid
		grid-template-columns:  1fr auto auto 1fr
		grid-template-rows: 1fr auto auto 1fr
		justify-content: center
		gap: var(--spacing)
		width: 100%
		
		.row
			grid-column: 1 / -1

		.property
			grid-column: 2
			display: flex
			align-items: end
		
		.value
			grid-column: 3
			font-size: 1.25rem
			display: flex
			align-items: end
		

	.remaining
		display: flex
		justify-content: right
		align-items: flex-end
		grid-row: 5

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

.error 
	color: var(--error-color)

.warn
	color: var(--warning-color)

.ok
	color: var(--success-color)

</style>
