<script lang="ts">
	import { AGORA_PROPOSAL_URL, EXPLORER_URL } from '@src/common/constants';

	import Button from '@src/components/starlight/components/Button.svelte';
	import BookOpenIcon from '@src/components/la/icons/book-open-solid.svelte';
	import SearchIcon from '@src/components/la/icons/search-solid.svelte';

	import type { CommonPeriodDetail } from '@src/common/types/governance';
	import { writeToClipboard } from '@src/util/clipboard';

	export let periodIndex: number;
	export let period: CommonPeriodDetail;

	function open_explorer() {
		window.open(`${EXPLORER_URL}/${period.proposal}`, '_blank');
	}

	function open_agora() {
		window.open(`${AGORA_PROPOSAL_URL}/${period.proposal}/${periodIndex}`, '_blank');
	}
</script>

<div
	class="period-wrap"
>
	<div class="period">
		<div class="proposal">
			<div class="title">
				<button class="unstyle-button" on:click={() => writeToClipboard(period.proposal)}>
					{period.proposal}
				</button>
			</div>

			<Button on:click={() => open_explorer()}>
				<div class="button-content">
					<div class="label">EXPLORER</div>
					<div class="icon"><SearchIcon /></div>
				</div>
			</Button>
			<Button on:click={() => open_agora()}>
				<div class="button-content">
					<div class="label">AGORA</div>
					<div class="icon"><BookOpenIcon /></div>
				</div>
			</Button>
		</div>
	</div>
</div>

<style lang="sass">
.period-wrap
	display: grid
	gap: var(--spacing)
	padding: var(--spacing)
	
	.period
		.proposal
			display: grid
			gap: var(--spacing)
			padding: var(--spacing)
			grid-template-columns: minmax(100px, 1fr) minmax(100px, 1fr) 

			.title
				display: flex
				justify-content: center
				font-size: 1.25rem
				font-weight: 500
				grid-column: 1 / -1	

				div
					min-width: 0
					white-space: nowrap
					text-overflow: ellipsis
					overflow: hidden

.button-content
	display: grid
	grid-template-columns: 1fr auto auto 1fr
	align-items: center
	gap: var(--spacing-f2)

	.label
		grid-column: 2

	.icon
		display: flex
		width: 20px
		fill: var(--button-text-color)
		align-items: center
</style>
