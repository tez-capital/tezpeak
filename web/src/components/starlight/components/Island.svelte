<script lang="ts">
	import AngleRight from '../icons/angle-right.svelte';
	import AngleLeft from '../icons/angle-left.svelte';
	import { ExpandingContainerDirection } from '../types';
	import Button from './Button.svelte';
	import Card from './Card.svelte';
	import ExpandingContainer from './ExpandingContainer.svelte';

	export let isOpen = false;
</script>

<div class="island-wrap">
	<Card>
		<div class="island">
			<div class="island-icon">
				<Button on:click={() => (isOpen = !isOpen)}>
					<div class="island-icon-content custom" class:open={isOpen}>
						<slot name="island-icon">
							<AngleLeft />
						</slot>
					</div>
					<div class="island-icon-content chevron" class:open={isOpen}>
						<AngleRight />
					</div>
				</Button>
			</div>

			<ExpandingContainer direction={ExpandingContainerDirection.horizontal} {isOpen}>
				<div class="island-content">
					<slot name="island-content" />
				</div>
			</ExpandingContainer>
		</div>
	</Card>
</div>

<style lang="sass">
.island-wrap
	--card-vertical-spacing: 0
	--card-horizontal-spacing: 0
.island
	display: grid
	grid-template-columns: 1fr auto
	grid-template-areas: "tools-menu tools-icon"
	
	--button-vertical-spacing: 0
	--button-horizontal-spacing: 0
	--button-height: var(--island-collapsed-height)
	--button-width: var(--island-collapsed-width)

	.island-icon
		cursor: pointer
		height: var(--island-collapsed-height)
		width: var(--island-collapsed-width)
		justify-content: center
		align-items: center
		display: flex
		grid-area: tools-icon

		.island-icon-content
			top: 0
			position: absolute
			display: flex
			justify-content: center
			align-items: center
			height: inherit
			width: inherit

		.chevron
			opacity: 0
			transform: scaleX(0)
			transition: all 0.2s ease-in-out

			&.open
				opacity: 1 
				transform: scaleX(1)

		.custom
			opacity: 1
			transition: all 0.2s ease-in-out

			&.open
				opacity: 0
				transform: scaleX(0)

		:global(svg)
			fill: var(--text-color)
			width: var(--island-icon-width)
			height: var(--island-icon-height)

</style>
