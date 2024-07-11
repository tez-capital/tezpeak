<script lang="ts">
	import AngleDownIcon from '../icons/angle-down.svelte';
	import AngleUpIcon from '../icons/angle-up.svelte';
	import { createKeyPressHandler, validate } from '../src/util';
	import type { ValidationRules } from '../types';
	import { v4 as uuidv4 } from 'uuid';

	export let id: string = uuidv4();
	export let value: number | string = ``;
	export let label: string = '';
	export let invalid: boolean = false;
	export let disabled: boolean = false;
	export let hint: string = '';
	export let type: 'text' | 'number' = 'text';
	export let rules: ValidationRules<any> = [];
	export let noLabel: boolean = false;

	let isInvalid = false;
	let invalidHint = '';

	let isFocused = false;
	function focus(focused: boolean) {
		isFocused = focused;
	}

	$: isInvalid = validate(value, rules) !== true;
	$: invalidHint =
		invalid === true
			? hint
			: (isInvalid ? validate(value, rules) as string : "");

	let inputElement: HTMLInputElement;
	function focusInput() {
		inputElement.focus();
	}
	function typeSetter(node: HTMLInputElement) {
		node.type = type;
	}

	// repeating for up/down arrows
	let interval: ReturnType<typeof setInterval>;
	let timeout: ReturnType<typeof setTimeout>;
	function startRepeatingAction(action: () => void) {
		action();

		timeout = setTimeout(() => {
			interval = setInterval(action, 50);
		}, 500);
	}

	function stopRepeatingAction() {
		clearTimeout(timeout);
		clearInterval(interval);
	}
</script>

<button
	class="unstyle-button input-wrap"
	class:input-wrap-hint-margin={hint !== undefined}
	class:focused={isFocused}
	class:disabled
	on:click={focusInput}
	on:keydown={(e) => ['Enter', ' '].includes(e.key) && focusInput()}
>
	<div class="input-wrap-internal" class:invalid={isInvalid}>
		<input
			bind:this={inputElement}
			{id}
			bind:value
			{...$$restProps}
			use:typeSetter
			on:focus={() => focus(true)}
			on:blur={() => focus(false)}
			class:shift={!noLabel}
		/>
		{#if !noLabel}
			<div class="label-wrap">
				<label
					for={id}
					class="label"
					class:label-shifted={isFocused || value.toString()}
					class:invalid-text={isInvalid}>{label}</label
				>
			</div>
		{/if}
		{#if type === 'number'}
			<div class="arrows">
				<button
					class="unstyle-button arrow"
					on:mousedown={() => startRepeatingAction(() => (value = Number(value) + 1))}
					on:mouseup={stopRepeatingAction}
					on:mouseleave={stopRepeatingAction}
					on:keydown={createKeyPressHandler(['Enter', ' '], () =>
						startRepeatingAction(() => (value = Number(value) + 1))
					)}
					on:keyup={createKeyPressHandler(['Enter', ' '], stopRepeatingAction)}
					on:focus={focusInput}
					aria-label="Decrease value"
				>
					<AngleUpIcon />
				</button>
				<button
					class="unstyle-button arrow"
					on:mousedown={() => startRepeatingAction(() => (value = Number(value) - 1))}
					on:mouseup={stopRepeatingAction}
					on:mouseleave={stopRepeatingAction}
					on:keydown={createKeyPressHandler(['Enter', ' '], () =>
						startRepeatingAction(() => (value = Number(value) - 1))
					)}
					on:keyup={createKeyPressHandler(['Enter', ' '], stopRepeatingAction)}
					on:focus={focusInput}
					aria-label="Decrease value"
				>
					<AngleDownIcon />
				</button>
			</div>
		{/if}
		{#if hint || invalidHint}
			<span
				class="hint"
				class:hidden={disabled || (!invalidHint && !hint)}
				class:invalid-text={isInvalid}
				class:hint-visible={isFocused || isInvalid}>{invalidHint || hint}</span
			>
		{/if}
	</div>
</button>

<style lang="sass">
.input-wrap-hint-margin
	margin-bottom: calc(var(--input-hint-font-size) + var(--input-hint-vertical-spacing) * 2) // hint + its padding

.input-wrap
	position: relative
	display: inline-block
	width: 100%
	border-radius: var(--input-border-radius)
	background-color: var(--input-background-color)
	cursor: text
	
	&.focused
		background-color: var(--input-background-highlight-color) 

	.input-wrap-internal
		position: relative
		display: grid
		grid-template-columns: minmax(50px, 1fr) auto
		grid-template-rows: var(--input-line-height)
		grid-template-areas: "text arrows"
		padding: var(--input-vertical-spacing)  var(--input-horizontal-spacing)
		border-radius: inherit
		box-sizing: border-box
		border: 1px solid transparent

		&.invalid
			border-color: var(--error-color) !important

	.label-wrap
		position: relative
		grid-area: text
		display: flex
		align-items: center
		justify-content: flex-start
		text-transform: none

	input
		grid-area: text
		font-family: inherit
		border: none
		border-radius: var(--border-radius)
		font-size: var(--font-size)
		outline: none
		color: var(--input-text-color)
		background: none
		width: calc(100% - var(--input-horizontal-spacing) * 2)
		margin: -1px -2px

		/* reset input */
		&:required, &:invalid
			box-shadow: none

		&::placeholder
			color: transparent

		&.shift
			transform: translateY(0.5rem)

	.hint
		display: block
		position: absolute
		left: 0
		top: 100%
		text-align: start
		width: calc(100% - var(--input-horizontal-spacing) * 2)
		height: 0
		color: var(--input-hint-color)
		font-size: var(--input-hint-font-size)
		opacity: 0
		transition-duration: 0.2s
		padding: var(--input-hint-vertical-spacing) var(--input-hint-horizontal-spacing)

	.hint-visible
		height: var(--input-hint-font-size)
		opacity: 1

	.invalid-text
		font-weight: bold
		color: var(--error-color) !important

	.label
		position: absolute
		left: 0
		display: block
		transition: 0.2s
		pointer-events: none
		color: var(--input-hint-color)
		overflow: hidden
		text-overflow: ellipsis
		white-space: nowrap
		text-transform: uppercase
	
	.label-shifted
		font-size: var(--input-label-font-size) !important
		font-weight: var(--input-label-font-weight)
		transform: translateY(-.8rem)


input[type=number]::-webkit-inner-spin-button
	appearance: none

.arrows
	--arrow-height: calc(var(--input-line-height) * .9)
	grid-area: arrows
	display: grid
	grid-template-columns: auto
	grid-template-rows: var(--arrow-height) var(--arrow-height)
	fill: var(--input-text-color)
	width: var(--arrow-height)
	align-self: center

	.arrow
		display: flex
		align-items: center !important
		height: 100%
		width: var(--arrow-height)
		
		&:hover
			transform: scale(1.1)
			cursor: pointer

</style>
