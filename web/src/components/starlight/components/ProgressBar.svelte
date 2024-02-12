<script lang="ts">
	import { onDestroy, onMount } from 'svelte';

	export let message = 'Loading...';
	export let hint = '';
	export let progress: number | 'indeterminate' = 0; // `indeterminate`;
	export let total = 100;

	let maskElement: HTMLDivElement;
	let progressElement: HTMLDivElement;

	$: percentage =
		progress === 'indeterminate' ? `-` : `${Math.min((progress / total) * 100, 100).toFixed(0)}%`;

	let observer: ResizeObserver;

	onMount(() => {
		observer = new ResizeObserver(() => {
			maskElement.style.width = `${progressElement.clientWidth}px`;
		});
		observer.observe(progressElement);
	});

	onDestroy(() => {
		observer.unobserve(progressElement);
	});
</script>

<div class="progressbar-wrap">
	<div class="data-grid message">
		<div class="data">{message}</div>
	</div>
	<div class="data-grid">
		<div class="data" class:opacity-0={progress === 'indeterminate'}>{percentage}</div>
	</div>
	<div bind:this={progressElement} class="progressbar-grid">
		<div class="line" class:identerminate={progress === 'indeterminate'}></div>
	</div>

	<div class="mask" class:opacity-0={progress === 'indeterminate'} style:width={percentage}>
		<div bind:this={maskElement} class="progressbar-grid" style:grid-row="1" style:grid-column="1">
			<div class="line"></div>
		</div>
	</div>
	<div class="data-grid">
		<div class="data hint">{hint}</div>
	</div>
</div>

<style lang="sass">
	.progressbar-wrap 
		position: relative
		display: grid
		grid-auto-rows: auto
		grid-template-columns: minmax(100px, 1fr)
		gap: var(--spacing-f2)
		width: var(--progress-bar-width, auto)
		height: var(--progress-bar-height, auto)

		--message-row: 1
		--progress-bar-row: 2
		--percentage-row: 3
	

	.data-grid 
		display: grid
		grid-template-columns: 1fr auto 1fr
		height: auto

		.data 
			grid-column: 2
			align-self: center
			justify-self: center
			pointer-events: none
		
	
	.message 
		padding-bottom: var(--spacing)
	

	.hint 
		font-size: var(--progress-bar-hint-font-size)
	

	.progressbar-grid 
		display: grid
		grid-row: var(--progress-bar-row)
		grid-template-columns: 1fr auto 1fr auto 1fr auto 1fr
	

	.line 
		align-self: center
		grid-column: 1 / 8
		grid-row: 1
		height: var(--progress-bar-line-height)
		background-color: var(--progress-bar-background-color)
	
	.identerminate 
		background: linear-gradient(90deg, transparent, var(--progress-bar-identerminate-color, white), transparent)
		background-size: 200% 100%
		animation: indeterminate-gradient 3s linear infinite
	

	.mask 
		grid-row: var(--progress-bar-row)
		overflow: hidden
		position: absolute
		left: 0
		width: 50%
		transition: width 0.3s

		:global(svg) 
			fill: var(--progress-bar-progress-color)
		

		.line 
			background-color: var(--progress-bar-progress-color)
		
	

	@keyframes indeterminate-gradient 
		0% 
			background-position: 100% 0%
		
		100% 
			background-position: -100% 0%
		
	
</style>
