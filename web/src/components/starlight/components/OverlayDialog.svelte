<svelte:options immutable />

<script lang="ts">
	let open: boolean = false;
	let persistent: boolean = false;
	let _class: string = '';
	let dialog: HTMLDialogElement;
	
	const close = () => {
		dialog.classList.add('hide');
		const cleanup = () => {
			dialog.classList.remove('hide');
			open = false;
			dialog.removeEventListener('animationend', cleanup);
		}
		dialog.addEventListener('animationend', cleanup);
	};
	$: background_click = persistent ? () => {} : close;

	export { _class as class, open, persistent };
</script>

<dialog class={_class} bind:this={dialog} {open}>
	<div class="content-wrap" aria-hidden="true" on:click|self={background_click}>
		<slot {close} />
	</div>
</dialog>

<style lang="sass">
dialog
	position: fixed
	top: 0
	left: 0
	width: 100vw
	height: 100vh
	padding: 0px
	margin: 0px 
	border: 0px
	background-color: var(--dialog-background-color)
	overflow: hidden
	z-index: 100

	.content-wrap
		display: flex
		align-items: center
		justify-content: center
		height: 100%
		width: 100%

dialog[open]
	animation: show-dialog .2s ease-in-out forwards

dialog.hide
	animation: hide-dialog .2s ease-in-out forwards

@keyframes show-dialog
	from
		opacity: 0
	to
		opacity: 1

@keyframes hide-dialog
	from
		opacity: 1
	to 
		opacity: 0


</style>
