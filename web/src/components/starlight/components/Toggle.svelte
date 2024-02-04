<svelte:options immutable />

<script lang="ts">
	import { createEventDispatcher } from 'svelte';

	export let label: string = '';
	export let checked: boolean = false;
	export let crossedOutIfNotChecked: boolean = false;
	let _class: string = '';
	export { _class as class };

	const dispatch = createEventDispatcher();
</script>

<label class="toggle-container {_class}">
	<div class="slider-container">
		<input type="checkbox" bind:checked on:change={() => dispatch('change', checked)} />
		<span class="slider"></span>
	</div>
	{#if label}
		<div class="label" class:crossed-out={crossedOutIfNotChecked && !checked}>{label}</div>
	{/if}
</label>

<style lang="sass">
.toggle-container 
	cursor: pointer
	display: inline-flex
	user-select: none
	align-items: center
	width: 100%

	.slider-container 
		position: relative
		border: var(--toggle-border)
		min-width: calc(var(--toggle-size) * 1.75)
		height: var(--toggle-size)
		
	.label 
		margin-left: var(--checkbox-label-padding-left)

		&.crossed-out 
			text-decoration: line-through

	.slider
		position: absolute
		cursor: pointer
		top: 0
		left: 0
		right: 0
		bottom: 0
		width: calc(var(--toggle-size) * 1.75)
		background-color: var(--toggle-background-color)
		border-radius: var(--toggle-border-radius)
		transition: .1s

	.slider:before 
		position: absolute
		content: ""
		height: calc(var(--toggle-size) * .8)
		width: calc(var(--toggle-size) * .8)
		left: calc(var(--toggle-size) * .1)
		bottom: calc(var(--toggle-size) * .1)
		background-color: var(--toggle-knob-color)
		transition: .1s
		border-radius: var(--toggle-border-radius)
	
	input
		opacity: 0

	input:checked + .slider
		background-color: var(--toggle-checked-background-color)
		

	input:focus + .slider
		box-shadow: 0 0 1px #2196F3
		

	input:checked + .slider:before
		transform: translateX(calc(var(--toggle-size) * .75))
		background-color: var(--toggle-checked-knob-color)


</style>
