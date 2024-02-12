<script lang="ts">
	import type { NavMenuItemType } from '../types';
	import NavMenuItem from './NavMenuItem.svelte';
	import NavHamburgerIcon from './NavHamburgerIcon.svelte';

	export let title: string = 'Menu';
	export let items: Array<NavMenuItemType> = [];
	export let bottomItems: Array<NavMenuItemType> = [];
	let expanded = false;
</script>

<div class="menu-layout" class:expanded>
	<div class="menu-items-wrap">
		<NavMenuItem text={title} on:click={() => (expanded = !expanded)}>
			<NavHamburgerIcon {expanded} slot="icon" />
		</NavMenuItem>
		<div class="gap" />
		<div class="items-list">
			{#each items as item}
				<!-- if not item.hidden -->
				{#if !item.hidden}
					<NavMenuItem
						text={item.label}
						icon={item.icon}
						isDisabled={item.disabled}
						path={item.path}
						on:click={item.click}
					/>
				{/if}
			{/each}
		</div>
	</div>
	<div class="items-list">
		{#each bottomItems as item}
			<!-- if not item.hidden -->
			{#if !item.hidden}
				<NavMenuItem
					text={item.label}
					icon={item.icon}
					isDisabled={item.disabled}
					path={item.path}
					on:click={item.click}
				/>
			{/if}
		{/each}
	</div>
</div>

<style lang="sass">
.menu-layout
	position: fixed
	display: grid
	grid-template-rows: 1fr auto
	height: calc(100vh - var(--spacing))
	width: var(--nav-menu-size)
	max-width: 100vw
	padding: var(--spacing-f2) 0
	
	background-image: var(--nav-menu-background-image)
	transition: 0.2s
	transition-timing-function: ease-in-out
	box-shadow: var(--nav-menu-box-shadow)
	overflow-x: hidden
	padding: var(--nav-menu-spacing)
	background-color: var(--nav-menu-background-color)
	border-radius: var(--nav-menu-border-radius)
	border-top-left-radius: var(--nav-menu-border-radius-tl, var(--nav-menu-border-radius, 0px))
	border-bottom-left-radius: var(--nav-menu-border-radius-bl, var(--nav-menu-border-radius, 0px))
	border-top-right-radius: var(--nav-menu-border-radius-tr, var(--nav-menu-border-radius, 0px))
	border-bottom-right-radius: var(--nav-menu-border-radius-br, var(--nav-menu-border-radius, 0px))

	.menu-items-wrap
		overflow: auto
		overflow-x: hidden

		.gap
			height: var(--nav-menu-spacing)

	.items-list
		display: grid
		grid-auto-rows: var(--nav-menu-item-height)
		overflow-x: hidden
		row-gap: var(--nav-menu-spacing)
		


	&.expanded
		width: var(--nav-menu-expanded-size)



</style>
