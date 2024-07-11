<script lang="ts">
	import { goto } from '$app/navigation';
	import Card from '@components/starlight/components/Card.svelte';
	import Button from '../starlight/components/Button.svelte';

	function open_governance() {
		goto('/tezpay');
	}
</script>

<div class="governance-wrap">
	<Card class="governance-card">
		<div class="governance">
			<div class="title">
				
			</div>
			<div/>
				<div class="payouts-info">
					<div class="kind">PAYOUTS</div>
					<div class="enter">
						<Button on:click={open_governance}>OPEN</Button>
					</div>
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
		grid-template-columns:  1fr
		grid-template-rows: 1fr auto auto auto 1fr
		justify-content: center
		gap: var(--spacing-f2)
		width: 100%
		
		.index, .kind
			display: flex
			justify-content: center
			font-size: 1.5rem
			font-weight: 500

		.index
			grid-row: 4

		.enter 
			display: flex
			justify-content: center
			grid-row: 3
			margin-top: var(--spacing-x2)
			margin-botton: var(--spacing-x2)

		.kind
			text-transform: uppercase
			grid-row: 2
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
</style>
