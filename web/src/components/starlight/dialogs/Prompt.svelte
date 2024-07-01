<script lang="ts">
	import OverlayDialog from '../components/OverlayDialog.svelte';
	import Card from '@components/starlight/components/Card.svelte';
	import type { PromiseFinalizers, ValidationRules } from '../types';
	import Input from '../components/Input.svelte';
	import { USER_CANCELED } from '../src/constants';
	import { validate } from '../src/util';
	import Button from '../components/Button.svelte';

	interface PromptDialogState<TValue> {
		title: string;
		message: string;
		hint: string;
		value: TValue;
		rules: ValidationRules;
		confirmText: string;
		cancelText: string;
	}

	const defaultState: PromptDialogState<any> = {
		title: 'Prompt',
		message: '',
		hint: '',
		value: '',
		rules: [],
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

	export async function prompt<TValue>(options: Partial<PromptDialogState<TValue>>): Promise<TValue> {
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

	$: isValid = validate(state.value, state.rules) !== true;
</script>

<OverlayDialog bind:open={isOpen} persistent>
	<Card>
		<div class="prompt-wrap">
			<slot name="title" {...state}>
				<h3 class="title">{state.title}</h3>
			</slot>
			<slot name="message" {...state}>
				{#if state.message}
					<p class="message">{state.message}</p>
				{/if}
			</slot>
			<slot name="value" {...state}>
				<Input {...state} bind:value={state.value} noLabel />
			</slot>
			<div class="controls padding-top">
				<div class="control-button" style:grid-column="1">
					<Button label={state.cancelText} on:click={() => promptFinalizers.reject(USER_CANCELED)} />
				</div>
				<div class="control-button" class:disabled={isValid} style:grid-column="3">
					<Button label={state.confirmText} on:click={() => promptFinalizers.resolve(state.value)} />
				</div>
			</div>
		</div>
	</Card>
</OverlayDialog>

<style lang="sass">
.prompt-wrap	
	position: relative
	display: grid
	grid-auto-rows: auto
	grid-template-columns: minmax(100px, 1fr)
	gap: var(--spacing)
	
	width: var(--prompt-dialog-width, auto)
	min-width: var(--prompt-dialog-min-width, 280px)
	max-width: var(--prompt-dialog-max-width, 500px)

	.title 
		text-align: center

	.message
		margin-top: 0px

	.controls
		display: grid
		grid-template-columns: auto 1fr auto

		.control-button
			min-width: 100px
</style>
