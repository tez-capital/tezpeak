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

<label class="checkbox-container {_class}">
	<div class="checkmark-container">
		<input type="checkbox" bind:checked on:change={() => dispatch('change', checked)} />
		<span class="checkmark"></span>
	</div>
	{#if label}
		<div class="label" class:crossed-out={crossedOutIfNotChecked && !checked}>{label}</div>
	{/if}
</label>

<style lang="sass">
.checkbox-container 
	cursor: pointer
	display: inline-flex
	user-select: none
	align-items: center
	width: 100%

	.checkmark-container 
		position: relative
		background-color: var(--checkbox-background-color)
		border-radius: var(--checkbox-border-radius)
		box-shadow: var(--checkbox-box-shadow)
		border: var(--checkbox-border)
		width: var(--checkbox-box-size)
		height: var(--checkbox-box-size)
		
	.label 
		margin-left: var(--spacing-f2)

		&.crossed-out 
			text-decoration: line-through

	input 
		width: var(--checkbox-box-size)
		height: var(--checkbox-box-size)
		opacity: 0
		cursor: pointer
		margin: 0

	input:checked ~ .checkmark 
		background-color: var(--checkbox-checked-background-color)
		box-shadow: var(--checkbox-checked-box-shadow)

	.checkmark:after 
		content: ""
		position: absolute
		opacity: 0
		transition: opacity 0.1s ease-in-out

	input:checked ~ .checkmark:after 
		opacity: 1

	.checkmark:after 
		left: 9px
		top: 5px
		width: 5px
		height: 10px
		border: solid var(--checkbox-checkmark-color)
		border-width: 0 3px 3px 0
		transform: rotate(45deg)
</style>
