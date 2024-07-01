<script lang="ts">
	import OverlayDialog from '../components/OverlayDialog.svelte';
	import Card from '@components/starlight/components/Card.svelte';
	import ProgressBar from '../components/ProgressBar.svelte';
	import type { ProgressValue } from '../types';

	interface ProgressOptions {
		progress: ProgressValue;
		message: string;
		hint: string;
		total: number;
	}

	interface ProgressDialogState extends ProgressOptions {
		title: string;
	}

	let isOpen = false;
	const defaultState: ProgressDialogState = {
		progress: 0,
		message: '',
		hint: '',
		total: 100,
		title: 'Progress'
	};
	let state = defaultState;

	export function show(options: Partial<ProgressDialogState>) {
		isOpen = true;
		state = { ...defaultState, ...options };
	}

	export function updateProgress(options: Partial<ProgressOptions>) {
		state = { ...state, ...options };
	}

	export function hide() {
		isOpen = false;
	}
</script>

<OverlayDialog bind:open={isOpen} persistent>
	<Card>
		<div class="progress-wrap">
			<div class="title">
				<h3>{state.title}</h3>
			</div>
			<slot name="progress-bar">
				<ProgressBar {...state} />
			</slot>
		</div>
	</Card>
</OverlayDialog>

<style lang="sass">
.progress-wrap
	width: var(--progress-dialog-width, auto)
	min-width: var(--progress-dialog-min-width, 280px)
	max-width: var(--progress-dialog-max-width, 500px)

	
	.title
		display: flex
		justify-content: center
		align-items: center

</style>
