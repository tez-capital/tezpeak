<svelte:options immutable />

<script lang="ts">
	export let value = false;
	export let label : string = '';

	const check = () => {
		value = !value;
	};
</script>

<div
	class="wrap"
	role="checkbox"
	aria-checked={value}
	tabindex="0"
	on:click|preventDefault={check}
	on:keydown={(e) => {
		if (e.key === 'Enter' || e.key === ' ') check();
	}}
>
	<label class="wrap-internal">
		<input type="checkbox" bind:checked={value} />
		<div class="checkmark-wrap">
			<div class="checkmark"></div>
		</div>
		{#if label}
			<div class="text">{label}</div>
		{/if}
	</label>
	
</div>

<style lang="sass">
.wrap 
	display: flex
	align-items: center
	cursor: pointer


.wrap-internal
	width: 100%
	height: auto
	cursor: pointer
	display: grid
	grid-template-columns: auto 1fr

	.checkmark-wrap
		grid-column: 1 / 2
		display: flex
		height: 100%
		align-items: center
		box-sizing: border-box

input[type='checkbox']
	display: none

.checkmark
	color: #ddd
	position: relative
	width: calc(var(--font-size) * 1.5)
	height: calc(var(--font-size) * 1.5)
	box-sizing: border-box

.text
	display: inline-flex
	margin-left: var(--spacing)
	text-align: left
	align-items: center

.checkmark:before
	content: ''
	display: inline-block
	position: absolute
	width: calc(var(--font-size) * 1.5)
	height: calc(var(--font-size) * 1.5)
	border: 1px solid var(--divider-color, white)
	opacity: 0
	left: 0
	top: 0
	box-sizing: border-box
	transition: all 0.12s, border-color 0.08s

.checkmark:after
	content: ''
	display: inline-block
	position: absolute
	width: calc(var(--font-size) * 1.5)
	height: calc(var(--font-size) * 1.5)
	border: 1px solid var(--divider-color, white)
	opacity: 0.6
	left: 0
	top: 0
	box-sizing: border-box
	transition: all 0.12s, border-color 0.08s

input[type='checkbox']:checked + .checkmark-wrap
	.checkmark:before
		width: calc(var(--font-size) * .85)
		opacity: 1
		left: 5px
		border-color: var(--checked-checkbox-color)
		border-top-color: transparent
		border-left-color: transparent
		margin-right: calc(var(--spacing-f2) + 6px)
		transform: rotate(45deg) translateY(-5px) translateX(0px)
</style>
