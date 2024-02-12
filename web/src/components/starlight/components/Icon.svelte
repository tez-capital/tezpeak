<svelte:options immutable />

<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { ComponentIcon } from '../types';

	export let icon: ComponentIcon;
	let _class: string = '';
	export { _class as class };

	const dispatch = createEventDispatcher();
</script>

<div class="icon {_class}" role="button" tabindex="0" on:click={() => dispatch('click')} on:keypress={() => dispatch('click')}>
	<svelte:component this={icon.component} />
</div>

<style lang="sass">
	.icon
		width: var(--icon-size)
		height: var(--icon-size)
		fill: var(--icon-color)
</style>
