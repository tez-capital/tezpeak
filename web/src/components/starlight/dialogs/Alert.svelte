<script lang="ts">
	import OverlayDialog from '../components/OverlayDialog.svelte';
	import Card from '@components/starlight/components/Card.svelte';
	import type { PromiseFinalizers, ValidationRules } from '../types';
	import Input from '../components/Input.svelte';
	import { USER_CANCELED } from '../src/constants';
	import { validate } from '../src/util';
	import Button from '../components/Button.svelte';

	interface AlertDialogState {
		title: string;
		message: string;
		hint: string;
	}

	const defaultState: AlertDialogState = {
		title: 'Prompt',
		message: '',
		hint: ''
	};
	export let state = defaultState;
	export let isOpen = false;

	let closeFinalizer: () => void;
	function close() {
		isOpen = false;
	}

	export async function alert(options: Partial<AlertDialogState>): Promise<void> {
		isOpen = true;
		state = { ...defaultState, ...options };
		return new Promise((resolve) => {
			closeFinalizer = () => {
				close();
				resolve();
			};
		});
	}
</script>

<OverlayDialog bind:open={isOpen}>
	<Card>
		<div class="alert-wrap">
			<slot name="title" {...state}>
				<h2 class="title">{state.title}</h2>
			</slot>
			<slot name="message" {...state}>
				{#if state.message}
					<p class="message">{state.message}</p>
				{/if}
			</slot>
			<div class="controls padding-top">
				<div class="control-button" style:grid-column="3">
					<Button label="close" on:click={() => closeFinalizer()} />
				</div>
			</div>
		</div>
	</Card>
</OverlayDialog>

<style lang="sass">
.alert-wrap	
	position: relative
	display: grid
	grid-auto-rows: auto
	grid-template-columns: minmax(100px, 1fr)
	gap: var(--spacing)
	
	width: var(--alert-dialog-width, auto)
	min-width: var(--alert-dialog-min-width, 280px)
	max-width: var(--alert-dialog-max-width, 500px)

	.title 
		text-align: center


	.message
		margin: var(--spacing-x2)
		margin-top: 0px
		margin-bottom: 0px

	.controls
		display: flex
		justify-content: center
		align-items: center

		.control-button
			min-width: 150px
</style>
