<script lang="ts">
	import { createEventDispatcher, onDestroy, onMount } from 'svelte';
	import { createMenuLayer, teleportBack, teleportToMenulayer } from '../src/util';
	import type { MenuItem } from '../types';
	import Button from './Button.svelte';
	import Card from './Card.svelte';

	export let align: 'left' | 'right' = 'right';
	export let items: Array<MenuItem> = [];

	let ownerElement: HTMLDivElement;
	let dialogElement: HTMLDivElement & { __origin?: HTMLDivElement | null };
	let itemsContainerElement: HTMLDivElement;

	let autoMenuWidth: string = 'unset';


	const dispatch = createEventDispatcher();

	const autoCloseListener = (e: MouseEvent) => {
		if (ownerElement.contains(e.target as Node)) return;
		close();
	};

	const close = () => {
		autoMenuWidth = 'unset';
		teleportBack(dialogElement);
		document.removeEventListener('click', autoCloseListener);
		dialogElement.style.removeProperty('--menu-width');
	};

	const open = async () => {
		const rect = ownerElement.getBoundingClientRect();

		dialogElement.style.top = `${rect.top}px`;
		let wantedWidth = 0;

		const style = getComputedStyle(dialogElement);
		switch (style.getPropertyValue('--menu-width')) {
			case '':
			case '"auto-inherit"':
				wantedWidth = Math.max(rect.width, itemsContainerElement.getBoundingClientRect().width);
				break;
			case '"inherit"':
				wantedWidth = rect.width;
				break;
			case 'auto':
				wantedWidth = itemsContainerElement.getBoundingClientRect().width;
				break;
			default:
				wantedWidth = dialogElement.getBoundingClientRect().width;
				break;
		}

		// position dialog
		let left = align === 'left' ? rect.left : rect.right - wantedWidth;
		if (left < 0) left = (window.innerWidth - wantedWidth) / 2;
		if (left + wantedWidth > window.innerWidth) left = (window.innerWidth - wantedWidth) / 2;

		dialogElement.style.left = `${left}px`;
		dialogElement.style.setProperty('--menu-width', `${wantedWidth}px`);
		document.addEventListener('click', autoCloseListener, { passive: true });
		teleportToMenulayer(dialogElement);
	};

	const select = (item: MenuItem) => {
		dispatch('select', item);
		close();
	};

	onMount(() => {
		createMenuLayer();
	});

	onDestroy(() => {
		close();
	});
</script>

<div class="container" bind:this={ownerElement}>
	<slot name="menu-button" {open}>
		<Button class="owner" on:click={open}>
			<slot name="menu-button-content">Menu</slot>
		</Button>
	</slot>
	<div class="dialog" bind:this={dialogElement}>
		<div class="items-container" bind:this={itemsContainerElement}>
			<Card class="items">
				{#each items as item}
					<slot name="menu-item" {select} {item}>
						<Button disabled="{item.disabled}" on:click={() => select(item)}>
							<div class="item">
								{#if item.icon}
									{#if item.icon.kind === 'font'}
										<i class="menu-item-icon {item.icon.family} {item.icon.id}" />
									{:else}
										<div class="menu-item-icon">
											<svelte:component this={item.icon.component} />
										</div>
									{/if}
								{/if}
								<span class="item-label">{item.label}</span>
							</div>
						</Button>
					</slot>
				{/each}
			</Card>
		</div>
	</div>
</div>

<style lang="sass">
.container
	display: inline-block
	width: 100%
	height: auto

:global(#sl-menu-layer)
	.dialog
		display: block !important

.dialog
	position: absolute
	display: none
	background: transparent
	margin: 0px
	padding: 0px
	border: 0px
	overflow-x: hidden
	min-width: var(--menu-min-width, 150px)
	width: var(--menu-width)
	max-height: var(--menu-max-height, calc(100vh - 2 * var(--spacing)))
	max-width: var(--manu-max-width, calc(100vw - 2 * var(--spacing)))

	--card-background-color: var(--menu-bacground-color)
	--card-vertical-spacing: var(--border-radius)
	--card-horizontal-spacing: 0px
	--button-border-radius: 0px
	--button-background-color: var(--card-background-color)
	--button-width: 100%
	--button-vertical-spacing: var(--spacing)

	:global(.items) 
		display: flex
		flex-direction: column

	.items-container
		display: inline-block
		width: var(--menu-width, auto)
		height: auto
		max-height: var(--menu-max-height, calc(100vh - 2 * var(--spacing) - 2 * var(--border-radius)))
		max-width: var(--manu-max-width, calc(100vw - 2 * var(--spacing)))

	.item
		position: relative
		display: inline-flex
		overflow: hidden
		align-items: center
		justify-content: var(--menu-justify-content, center)
		width: 100%

		.item-label
			overflow: hidden
			text-overflow: ellipsis
			text-transform: none
			white-space: nowrap
								
		.menu-item-icon
			margin-right: var(--spacing)
			color: var(--text-color)
			fill: var(--text-color)
			width: var(--menu-item-icon-size, 1.5rem)
			min-width: var(--menu-item-icon-size, 1.5rem)
			height: var(--menu-item-icon-size, 1.5rem)
			min-height: var(--menu-item-icon-size, 1.5rem)


</style>
