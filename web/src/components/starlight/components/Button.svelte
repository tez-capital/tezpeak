<svelte:options immutable />

<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	export let label: string = '';
	export let text: string = '';
	export let disabled: boolean = false;
	let _class: string = '';
	export { _class as class };

	$: actualText = text || label;

	const dispatch = createEventDispatcher();
</script>

<button class="button {_class}" class:disabled={disabled} on:click={() => dispatch('click')}>
	<slot>
		{actualText}
	</slot>
</button>

<style lang="sass">
.button 
	cursor: pointer
	display: inline-block
	background-color: var(--button-background-color)
	background-image: var(--button-background-image)
	color: var(--button-text-color)
	padding: var(--button-vertical-spacing)  var(--button-horizontal-spacing)
	border-radius: var(--button-border-radius)
	font-size: var(--button-font-size)
	font-weight: var(--button-font-weight)
	white-space: var(--button-whitespace)
	text-transform: var(--button-text-transform)
	text-align: var(--button-text-align)
	text-decoration: var(--button-text-decoration)
	border: var(--button-border)
	outline: var(--button-outline)
	font-family: var(--button-font-family)
	overflow: var(--button-overflow)
	text-overflow: var(--button-text-overflow)
	letter-spacing: var(--button-letter-spacing)
	transition: var(--button-transition)
	background-size: var(--button-background-size)
	box-shadow: var(--button-box-shadow)
	background-position: var(--button-background-position)
	width: var(--button-width)
	height: var(--button-height)

	&:hover 
		color: var(--button-hover-text-color)
		background-color: var(--button-hover-background-color)
		box-shadow: var(--button-hover-box-shadow) !important

	&:active
		box-shadow: var(--button-active-box-shadow) !important

	svg
		fill: var(--button-text-color)
		stroke: var(--button-text-color)

</style>
