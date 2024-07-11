<script lang="ts">
	import { ExpandingContainerDirection } from '../types';

	export let isOpen = false;
	export let direction: ExpandingContainerDirection = ExpandingContainerDirection.vertical;

	let contentElement: HTMLDivElement;
	let containerElement: HTMLDivElement;

	function computeHeight(element: HTMLDivElement) {
		if (!element) {
			return ``;
		}
		const computedStyle = getComputedStyle(element);
		const height = computedStyle.getPropertyValue('--expanding-container-expanded-height');
		if (height === 'auto') {
			return `${contentElement.clientHeight}px`;
		}
		return height;
	}

	function computeWidth(element: HTMLDivElement) {
		if (!element) {
			return ``;
		}
		const computedStyle = getComputedStyle(element);
		const width = computedStyle.getPropertyValue('--expanding-container-expanded-width');
		if (width === 'auto') {
			return `${contentElement.clientWidth}px`;
		}
		return width;
	}

	$: contentHeight = isOpen ? computeHeight(containerElement) : `0px`;
	$: contentWidth = isOpen ? computeWidth(containerElement) : `0px`;
</script>

<div
	class:container-vertical={direction === ExpandingContainerDirection.vertical}
	class:container-horizontal={direction === ExpandingContainerDirection.horizontal}
	bind:this={containerElement}
	style:height={direction === ExpandingContainerDirection.vertical ? contentHeight : ''}
	style:width={direction === ExpandingContainerDirection.horizontal ? contentWidth : ''}
	class:open={isOpen}
	class:overflow-auto={contentHeight === ''}
>
	<div bind:this={contentElement} class="content">
		<slot />
	</div>
</div>

<style lang="sass">
	.container-vertical 
		height: var(--expanding-container-collapsed-height, 0)
		max-height: var(--expanding-container-max-height)
		overflow: var(--expanding-container-overflow, hidden)
		transition: height var(--expanding-container-transition-duration) ease-in-out
		position: relative

		&.open 
			height: var(--expanding-container-expanded-height, auto)

		
		.content 
			height: auto
			padding: var(--expanding-container-padding, 0)

	.container-horizontal 
		width: var(--expanding-container-collapsed-width, 0)
		max-width: var(--expanding-container-max-width)
		overflow: var(--expanding-container-overflow, hidden)
		transition: width var(--expanding-container-transition-duration) ease-in-out
		position: relative
		height: 100%

		&.open 
			width: var(--expanding-container-expanded-width, auto)
	
		.content 
			height: 100%
			width: max-content
			padding: var(--expanding-container-padding, 0)
</style>
