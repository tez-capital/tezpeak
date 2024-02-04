<script lang="ts">
	import AngleDownIcon from '../icons/angle-down.svelte';
	import { createEventDispatcher, onMount } from 'svelte';
	import type { SelectItem } from '../types';
	import Menu from './Menu.svelte';

	export let label: string = 'Options';
	export let options: SelectItem[] = [];
	export let value: SelectItem | undefined = undefined;
	export let selectedItemIndex = -1;

	const dispatch = createEventDispatcher();

	const select = (item: SelectItem) => {
		value = item;
		dispatch('selectionChanged', value);
	};
	onMount(() => {
		if (selectedItemIndex >= 0) {
			value = options[selectedItemIndex];
		}
	});
</script>

<div class="container">
	<Menu items={options} on:select={(event) => select(event.detail)}>
		<div class="select" slot="menu-button-content">
			<div class="text">
				<span class="label" class:selected={value}>{label}</span>
				{#if value}
					<span class="value">{value?.label}</span>
				{/if}
			</div>

			<span class="arrow">
				<AngleDownIcon />
			</span>
		</div>
	</Menu>
</div>

<style lang="sass">
.container
	position: relative
	width: 100%

	.select
		display: grid
		grid-template-columns: minmax(0, 1fr) auto
		grid-template-areas: "text arrow"

		.text
			position: relative
			grid-area: text
			display: flex
			align-items: center
			justify-content: flex-start
			text-transform: none

			span
				overflow: hidden
				text-overflow: ellipsis
				white-space: nowrap
		
		.label
			position: absolute
			left: 0
			color: var(--text-color)
			text-transform: uppercase
		
			&.selected
				font-size: calc(var(--font-size) * .6) 
				transform: translateY(-.8rem)
				color: var(--hint-color)
				font-weight: var(--select-label-font-weight)
		.value
			transform: translateY(0.5rem)

		.arrow
			grid-area: arrow
			color: var(--text-color)
			fill: var(--text-color)
			width: var(--select-item-icon-size, 1.5rem)
			min-width: var(--select-item-icon-size, 1.5rem)
			height: var(--select-item-icon-size, 1.5rem)
			min-height: var(--select-item-icon-size, 1.5rem)

</style>
