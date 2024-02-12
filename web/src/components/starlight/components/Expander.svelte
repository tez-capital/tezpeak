<script lang="ts">
	import ChevronDownSolid from '@src/components/la/icons/chevron-down-solid.svelte';
	import Card from './Card.svelte';
	import ExpandingContainer from './ExpandingContainer.svelte';
	import ArrowDownSolid from '@src/components/la/icons/arrow-down-solid.svelte';
	import AngleDown from '../svg/angle-down.svelte';
	import Button from './Button.svelte';

	export let isOpen = false;
	export let title = 'Expander';
</script>

<div class="expander-wrap">
	<Card>
		<div class="expander-grid">
			<slot name="header">
				<Button on:click={() => (isOpen = !isOpen)}>
					<div class="expander-header">
						<slot name="header-title">
							{title}
						</slot>
						<div class="icon-wrap">
							<slot name="icon">
								<div class="icon" class:open={isOpen}>
									<AngleDown />
								</div>
							</slot>
						</div>
					</div>
				</Button>
			</slot>

			<ExpandingContainer {isOpen}>
				<div class="expander-content">
					<slot />
				</div>
			</ExpandingContainer>
		</div>
	</Card>
</div>

<style lang="sass">
.expander-wrap
	--card-vertical-spacing: 0
	--card-horizontal-spacing: 0

.expander-grid
	display: grid
	grid-template-columns: 1fr
	grid-template-rows: auto auto
	grid-template-areas: "header" "content"

	.expander-header
		grid-area: header
		cursor: pointer
		display: grid
		grid-template-columns: auto 1fr auto
		grid-template-areas: "title _ icon"
		width: 100%
		color: var(--text-color)
		
		.icon-wrap
			grid-area: icon
			display: flex
			justify-content: center
			align-items: center
			height: 100%

			.icon
			
				fill: var(--text-color)
				width: var(--font-size)
				
				transition: transform 0.2s ease-in-out

				&.open
					transform: scaleY(-1)



	.expander-content
		grid-area: content
		padding: var(--spacing)
</style>
