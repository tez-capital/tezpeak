<script lang="ts">
	import { Button, Card } from '@alis-is/starlight';

	interface Props {
		ontestnotifications: () => void;
		ontestextensions: () => void;
	}
	let { ontestnotifications, ontestextensions }: Props = $props();
</script>

<Card>
	{#snippet header()}
		<h3 class="sl-card-title">Test</h3>
	{/snippet}

	<p class="hint">Send a test notification or run the configured extensions without paying anything.</p>

	{#snippet footer()}
		<Button onclick={ontestnotifications}>Notifications</Button>
		<Button onclick={ontestextensions}>Extensions</Button>
	{/snippet}
</Card>

<style lang="sass">
p
	margin: 0
</style>
