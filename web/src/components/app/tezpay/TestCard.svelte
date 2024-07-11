<script lang="ts">
	import Separator from '../Separator.svelte';
	import Card from '@components/starlight/components/Card.svelte';
	import { createEventDispatcher, onDestroy } from 'svelte';
	import Button from '../../starlight/components/Button.svelte';

	const dispatch = createEventDispatcher<{
		preview: void;
		pay: void;
	}>();
</script>

<div class="test-wrap">
	<Card class="test-card">
		<div class="test">
			<div class="title">
				<h5>TEST</h5>
			</div>
			<Separator />

			<div class="tools">
				<Button on:click={() => dispatch('preview')}>Notifications</Button>
				<Button on:click={() => dispatch('pay')}>Extensions</Button>
			</div>
		</div>
	</Card>
</div>

<style lang="sass">
.test-wrap
	display: grid
	grid-template-rows: 1fr
	width: 100%
	height: 100%
	user-select: none
	// &:hover
	// 	cursor: pointer
	// 	transition: background-color 0.2s
	// 	--card-background-color: #151515

	:global(.test-card)
		box-sizing: border-box
		height: 100%

.test
	display: grid
	grid-template-rows: auto auto auto auto 1fr
	height: 100%
	gap: var(--spacing)
	.title
		display: flex
		justify-content: center
		h5
			font-size: 1.5rem
			font-weight: 500
			margin: 0
	.tools
		display: flex
		flex-direction: column
		justify-content: center
		align-items: center
		height: 100%
		gap: var(--spacing)


</style>
