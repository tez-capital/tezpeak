<script lang="ts">
	import Card from '@components/starlight/components/Card.svelte';
	import Separator from '@components/peak/Separator.svelte';
	import { writeToClipboard } from '@src/util/clipboard';
	import { formatBalance } from '@src/util/format';
	import { calculateFreeSpace, calculateStakingCapacity, getBakerColor } from '@src/util/baker';
	import type { BakerStatus } from '@src/common/types/status';

	export let baker: string;
	export let status: BakerStatus;
	export let showColor = false;

	$: bakerColor = getBakerColor(baker);

</script>

<Card>
	<div class="baker-grid">
		<div class="title">
			<h5>
				Baker
				{#if showColor}
					<div class="assigned-color" style:--assigned-color-color={bakerColor}>
						<div class="assigned-color-sign"></div>
					</div>
				{/if}
			</h5>
			<button class="unstyle-button address" on:click={() => writeToClipboard(baker)}
				>{baker}</button
			>
		</div>
		<Separator />
		<div class="balances">
			<div class="title">Balance</div>
			<div class="value">{formatBalance(status.full_balance)}</div>
			<div class="title">Frozen deposit limit</div>
			<div class="value">{formatBalance(status.frozen_deposits_limit)}</div>
			<div class="title">Staking</div>
			<div class="value">{formatBalance(status.staking_balance)}</div>
		</div>
		<Separator />
		<div class="balances">
			<div class="title">Staking Capacity</div>
			<div class="value">{formatBalance(calculateStakingCapacity(status))}</div>
			<div class="title">Free Space</div>
			<div class="value">{formatBalance(calculateFreeSpace(status))}</div>
		</div>
		<div class="bottom-separator">
			<Separator />
		</div>
		<div class="delegators">
			<div class="value">{status.delegated_contracts.length - 1}</div>
			delegators
		</div>
	</div>
</Card>

<style lang="sass">
	.baker-grid
		display: grid
		grid-template-columns: minmax(100px, 1fr)
		gap: var(--spacing)

		.title
			position: relative

			h5
				font-size: 1.5rem
				font-weight: 500
				margin: 0
			
			.address
				margin: 0
				margin-top: var(--spacing-f2)
				font-size: 0.95rem
		.balances
			display: grid
			grid-template-columns: auto 1fr auto
			gap : var(--spacing-f2)
			
			div
				display: flex
				white-space: nowrap
				align-items: center
				margin: var(--spacing-f2) 0
			
			.title
				grid-column: 1
				justify-content: left

			.value
				grid-column: 3
				justify-content: right

		.bottom-separator
			grid-row: 6
			display: flex
			align-items: flex-end
		
		.delegators
			grid-row: 7
			.value
				font-size: 1.25rem
				font-weight: 500
				display: inline-block

.assigned-color
	.assigned-color-sign
		position: absolute
		display: inline-block
		width: 12px
		height: 12px
		border-radius: 20%
		margin-right: var(--spacing-f2)
		background-color: var(--assigned-color-color)
		top: 0
		right: 0
</style>
