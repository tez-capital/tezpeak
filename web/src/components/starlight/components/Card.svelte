<svelte:options immutable />

<script lang="ts">
	let _class: string = '';
	export { _class as class };
</script>

<div class="card {_class}">
	<slot/>
</div>

<style lang="sass">
.card
	position: relative
	color: var(--card-text-color)
	background-color: var(--card-background-color)
	background-image: var(--card-background-image)
	padding: var(--card-vertical-spacing)  var(--card-horizontal-spacing)
	border-radius: var(--card-border-radius)
	border: var(--card-border)
	outline: var(--card-outline)
	overflow: var(--card-overflow)
	background-size: var(--card-background-size)
	box-shadow: var(--card-box-shadow)
	background-position: var(--card-background-position)
	width: var(--card-width)
	height: var(--card-height)
	max-width: var(--card-max-width)
	max-height: var(--card-max-height)
	overflow: var(--card-overflow)

</style>
