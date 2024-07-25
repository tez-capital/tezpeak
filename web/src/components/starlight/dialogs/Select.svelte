<script lang="ts">
	import OverlayDialog from './Overlay.svelte';
	import Card from '@components/starlight/components/Card.svelte';
	import type { PromiseFinalizers, SelectItem } from '../types';
	import { USER_CANCELED } from '../src/constants';
	import Button from '../components/Button.svelte';
	import Select from '../components/Select.svelte';

	interface SelectDialogState<TValue> {
		title: string;
		message: string;
		hint: string;
		value: TValue;
		options: Array<SelectItem<TValue> | TValue>;
		confirmText: string;
		cancelText: string;
	}

	const defaultState: SelectDialogState<any> = {
		title: 'Select',
		message: '',
		hint: '',
		value: '',
		options: [],
		confirmText: 'Confirm',
		cancelText: 'Cancel'
	};
	export let state = defaultState;
	export let isOpen = false;
	let value: SelectItem<any> | undefined = undefined;
	let options: Array<SelectItem<any>> = [];

	let promptFinalizers: PromiseFinalizers = {
		resolve: () => close(),
		reject: () => close()
	};

	function close() {
		isOpen = false;
	}

	export async function prompt<TValue>(
		promptOptions: Partial<SelectDialogState<TValue>>
	): Promise<TValue> {
		isOpen = true;
		options = (promptOptions.options ?? []).map(
			(v) => {
				if (typeof v === 'object') {
					return v as SelectItem<TValue>;
				}
				return { label: v, value: v } as SelectItem<TValue>;
			}
		);

		state = { ...defaultState, ...promptOptions };
		value = options.find((v) => (v as SelectItem<TValue>).value === state.value) as SelectItem<TValue>;
		
		return new Promise<TValue>((resolve, reject) => {
			promptFinalizers = {
				resolve: () => {
					close();
					resolve(value?.value);
				},
				reject: (e) => {
					close();
					reject(e);
				}
			};
		});
	}

	$: isValid = state.options.includes(value);
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
			<Select bind:value options={options} />
			<div class="controls padding-top">
				<div class="control-button" style:grid-column="1">
					<Button
						label={state.cancelText}
						on:click={() => promptFinalizers.reject(USER_CANCELED)}
					/>
				</div>
				<div class="control-button" class:disabled={isValid} style:grid-column="3">
					<Button
						label={state.confirmText}
						on:click={() => promptFinalizers.resolve(state.value)}
					/>
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
	
	width: var(--select-dialog-width, auto)
	min-width: var(--select-dialog-min-width, 280px)
	max-width: var(--select-dialog-max-width, 500px)

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
