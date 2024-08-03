<script lang="ts">
	import OverlayDialog from './Overlay.svelte';
	import Card from '@components/starlight/components/Card.svelte';
	import type { PromiseFinalizers } from '../types';
	import { USER_CANCELED } from '../src/constants';
	import Button from '../components/Button.svelte';

	interface ConfirmDialogState {
		title: string;
		message: string;
		confirmText: string;
		cancelText: string;
	}

	const defaultState: ConfirmDialogState = {
		title: 'Confirm',
		message: '',
		confirmText: 'Confirm',
		cancelText: 'Cancel'
	};
	export let state = defaultState;
	export let isOpen = false;

	let promptFinalizers: PromiseFinalizers = {
		resolve: () => close(),
		reject: () => close()
	};

	function close() {
		isOpen = false;
	}

	export async function request_confirmation<TValue>(options: Partial<ConfirmDialogState>): Promise<TValue> {
		isOpen = true;
		state = { ...defaultState, ...options };
		return new Promise<TValue>((resolve, reject) => {
			promptFinalizers = {
				resolve: (v) => {
					close();
					resolve(v);
				},
				reject: (e) => {
					close();
					reject(e);
				}
			};
		});
	}
</script>

<OverlayDialog bind:open={isOpen} persistent>
	<Card>
		<div class="confirm-wrap">
			<slot name="title" {...state}>
				<h3 class="title">{state.title}</h3>
			</slot>
			<slot name="message" {...state}>
				{#if state.message}
					<p class="message">{state.message}</p>
				{/if}
			</slot>
			<div class="controls padding-top">
				<div class="control-button" style:grid-column="1">
					<Button label={state.cancelText} on:click={() => promptFinalizers.reject(USER_CANCELED)} />
				</div>
				<div class="control-button" style:grid-column="3">
					<Button label={state.confirmText} on:click={() => promptFinalizers.resolve()} />
				</div>
			</div>
		</div>
	</Card>
</OverlayDialog>

<style lang="sass">
.confirm-wrap	
	position: relative
	display: grid
	grid-auto-rows: auto
	grid-template-columns: minmax(100px, 1fr)
	gap: var(--spacing)
	
	width: var(--confirm-dialog-width, auto)
	min-width: var(--confirm-dialog-min-width, 280px)
	max-width: var(--confirm-dialog-max-width, 500px)
	padding: var(--confirm-dialog-padding, var(--spacing-x2))

	.title 
		text-align: center

	.message
		margin-top: 0px
		padding: var(--spacing)

	.controls
		display: grid
		grid-template-columns: auto 1fr auto

		.control-button
			min-width: 100px
			--button-horizontal-spacing: var(--spacing-x2)
			
</style>
