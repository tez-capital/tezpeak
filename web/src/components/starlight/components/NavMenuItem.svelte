<svelte:options immutable />

<script lang="ts">
	import { createEventDispatcher } from 'svelte';
	import type { IIcon } from '../types';
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';

	export let text = '';
	export let icon: IIcon = { kind: 'font', id: 'las-home', family: 'la' };
	export let isDisabled = false;
	export let isActive: boolean | undefined = undefined;
	export let path: string | null = null;

	const dispatch = createEventDispatcher();

	const navigate = () => {
		dispatch("click");
		if (path && path !== '.') goto(path);
	};

	$: active = isActive !== undefined ? isActive : path && $page.url.pathname === path;
	
</script>

<div class="menu-item-wrap" class:disabled={isDisabled} class:active={active}>
	<a class="menu-item" href={path} on:click|preventDefault={navigate}>
		<slot name="icon">
			<div class="menu-item-icon svg-icon">
				{#if icon.kind === 'font'}
					<i class="{icon.family} {icon.id}" />
				{:else if icon.kind === 'component'}
					<svelte:component this={icon.component} />
				{/if}
			</div>
		</slot>
		<div class="menu-item-text">
			<span>{text}</span>
		</div>
	</a>
</div>

<style lang="sass">
.menu-item-wrap
	overflow: hidden
	position: relative
	color: var(--nav-menu-text-color)
	
	.menu-item 
		display: grid
		grid-template-columns: var(--nav-menu-size) minmax(40px, 1fr)
		grid-template-areas: "icon text"
		height: var(--nav-menu-item-height)
		border-radius: var(--nav-menu-item-border-radius)
		width: 100%
		transition: var(--nav-menu-item-transition)
		cursor: pointer
		text-decoration: none
		background-color: var(--nav-menu-item-background-color)
		
		&:hover
			background-color: var(--nav-menu-item-hover-background-color)

		.menu-item-icon
			grid-area: icon
			display: flex
			justify-content: center
			align-items: center
			position: relative
	
			--svg-icon-width: var(--nav-menu-icon-size)
			--svg-icon-height: var(--nav-menu-icon-size)
			--svg-icon-fill: var(--nav-menu-text-color)

			i
				font-size: var(--nav-menu-icon-size)


		.menu-item-text
			grid-area: text
			display: flex
			justify-content: flex-start
			align-items: center
			color: var(--nav-menu-text-color)
			margin-right: var(--spacing)

			span
				overflow: hidden
				text-overflow: ellipsis
				text-align: left
				white-space: nowrap
				font-size: var(--nav-menu-font-size)

	.disabled
		pointer-events: none
		filter: grayscale(100%) opacity(50%)

.active
	--nav-menu-item-background-color: var(--nav-menu-item-active-background-color, --nav-menu-item-background-color)
	--nav-menu-text-color: var(--nav-menu-item-active-text-color, --nav-menu-text-color)
	--nav-menu-item-hover-background-color: var(--nav-menu-item-active-hover-background-color, --nav-menu-item-hover-background-color)

	&:hover
		--nav-menu-item-active-background-color: var(--nav-menu-item-active-hover-background-color, --nav-menu-item-active-background-color)

	span
		font-weight: var(--nav-menu-item-active-font-weight, bold)
</style>
