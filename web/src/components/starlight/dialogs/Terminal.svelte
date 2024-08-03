<script lang="ts">
	import OverlayDialog from './Overlay.svelte';
	import Card from '@components/starlight/components/Card.svelte';
	import Button from '../components/Button.svelte';
	import Terminal from '../components/Terminal.svelte';

	interface TerminalDialogState {
		title: string;
		allowClose: boolean;
		closeText: string;
	}

	const defaultState: TerminalDialogState = {
		title: 'Select',
		allowClose: true,
		closeText: 'Cancel'
	};

	let terminal: Terminal;
	export let state = defaultState;
	export let isOpen = false;

	let closeFinalizer: () => void;

	function close() {
		isOpen = false;
	}

	export async function hide() {
		close();
	}

	export async function write(data: string) {
		terminal.write(data);
	}

	export async function clear(data: string) {
		terminal.clear();
	}

	export async function update_state(options: Partial<TerminalDialogState>) {
		state = { ...state, ...options };
	}

	export async function show(dialogOptions: Partial<TerminalDialogState>): Promise<void> {
		isOpen = true;
		state = { ...defaultState, ...dialogOptions };

		return new Promise((resolve) => {
			closeFinalizer = () => {
				close();
				resolve();
			};
		});
	}
</script>

<OverlayDialog bind:open={isOpen} persistent>

	<Card>
		<div class="terminal-wrap">
			<slot name="title" {...state}>
				<h3 class="title">{state.title}</h3>
			</slot>
			<div class="terminal">
				<Terminal bind:this={terminal} />
			</div>
			<div class="controls padding-top">
				<div class="control-button" class:disabled={!state.allowClose} style:grid-column="2">
					<Button label={state.closeText} on:click={() => closeFinalizer()} />
				</div>
			</div>
		</div>
	</Card>
</OverlayDialog>

<style lang="sass">
.terminal-wrap	
	position: relative
	display: grid
	grid-auto-rows: auto 1fr auto
	grid-template-columns: minmax(100px, 1fr)
	gap: var(--spacing)
	
	height: var(--terminal-dialog-height, 80vh)
	min-height: var(--terminal-dialog-min-height, 280px)
	max-height: var(--terminal-dialog-max-height)
	width: var(--terminal-dialog-width, 80vw)
	min-width: var(--terminal-dialog-min-width, 280px)
	max-width: var(--terminal-dialog-max-width)

	.title 
		text-align: center

	.terminal
		margin: var(--spacing)

	.controls
		display: grid
		grid-template-columns: auto 1fr auto

		.control-button
			min-width: 150px
			--button-horizontal-spacing: var(--spacing-x2)
</style>
