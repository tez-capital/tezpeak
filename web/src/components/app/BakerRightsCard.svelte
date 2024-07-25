<script lang="ts">
	import { flip } from 'svelte/animate';
	import type { BlockRights } from '@src/common/types/status';
	import Card from '../starlight/components/Card.svelte';
	import Separator from '@components/app/Separator.svelte';
	import { getBakerColor, normalizeBlockRights } from '@src/util/baker';

	import BlockIcon from '@components/la/icons/cube-solid.svelte';
	import AttestationIcon from '@components/la/icons/shield-alt-solid.svelte';

	export let mode: 'upcoming' | 'past';
	export let title = 'Baking Rights';
	export let rights: BlockRights[];
	export let showBakerColors = false;
</script>

<Card>
	<div class="baker-rights">
		<div class="title">
			<h5>{title}</h5>
		</div>
		<Separator />
		<div class="block-rights-wrap">
			{#each rights as blockRights (blockRights.level)}
				<div class="block-rights" animate:flip>
					<div class="level">
						{blockRights.level}
					</div>

					<div class="block-rights-inner">
						{#each normalizeBlockRights(blockRights) as right}
							{#if showBakerColors}
								<div
									class="baker-assigned-color"
									style:--assigned-color-color={getBakerColor(right.baker)}
								>
									<div class="baker-assigned-color-sign"></div>
								</div>
							{/if}
							<div class="baker-block-rights">
								<div
									class="value"
									class:no-rights={right.blocks === 0}
									class:warning={mode === 'past' && right.realizedBlocks === 0}
									class:success={mode === 'past' && right.realizedBlocks === 1}
								>
									{#if mode === 'past'}
										{right.blocks}/{right.realizedBlocks * right.blocks}
									{:else}
										{right.blocks}
									{/if}
								</div>
								<div class="icon" class:no-rights={right.blocks === 0}>
									<BlockIcon />
								</div>

								<div
									class="value"
									class:no-rights={right.attestations === 0}
									class:warning={mode === 'past' && right.realizedAttestations === 0}
									class:success={mode === 'past' && right.realizedAttestations === 1}
								>
									{#if mode === 'past'}
										{right.attestations}/{right.realizedAttestations * right.attestations}
									{:else}
										{right.attestations}
									{/if}
								</div>
								<div class="icon" class:no-rights={right.attestations === 0}>
									<AttestationIcon />
								</div>
							</div>
						{/each}
					</div>
				</div>
			{/each}
		</div>
	</div>
</Card>

<style lang="sass">
.baker-rights
	display: grid
	gap: var(--spacing)
	.title
		display: flex
		justify-content: center
		h5
			font-size: 1.5rem
			font-weight: 500
			margin: 0
	.block-rights-wrap
		display: grid

		max-height: 400px
		overflow: auto
		justify-content: center

		.block-rights
			display: grid
			gap: var(--spacing-f2)
			grid-template-columns: auto auto 1fr
			align-items: center
			transition: background-color 0.2s
			padding: var(--spacing-f2)
			border-radius: var(--border-radius)

			&:hover
				background-color: rgba(255,255,255, 0.15)
				cursor: pointer
			
			.level
				grid-column: 1

			.block-rights-inner
				grid-column: 2
				display: flex
				align-items: center
				min-width: 150px
			
				
				.no-rights
					filter: grayscale(0.75) opacity(0.44)
				.baker-block-rights
					// display: inline-flex
					// align-items: center
					display: grid
					grid-template-columns: 120px 30px 120px 30px 

					.value
						display: inline-flex
						justify-content: right
						align-items: center
						white-space: nowrap
						margin-right: var(--spacing-f2)
						&.warning
							color: var(--warning-color)
						
						&.success
							color: var(--success-color)

	.icon
		display: inline-block
		width: 30px
		height: 30px
		fill: var(--text-color)

.baker-assigned-color
	display: inline-block

	.baker-assigned-color-sign
		display: inline-block
		width: 20px
		height: 20px
		border-radius: 20%
		margin-right: var(--spacing-f2)
		background-color: var(--assigned-color-color)
</style>
