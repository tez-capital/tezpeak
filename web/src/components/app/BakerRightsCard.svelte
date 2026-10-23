<script lang="ts">
	import { goto } from '$app/navigation';
	import { Card, Button } from '@alis-is/starlight';
	import { state as tezbakeState } from '@app/state/tezbake';
	import { bakerRightsTimeline, getBakerColor, summarizeRights } from '@src/util/baker';
	import { formatAddress } from '@src/util/format';
	import RightsStrip from './RightsStrip.svelte';

	interface Props {
		baker: string;
		showColor?: boolean;
	}
	let { baker, showColor = false }: Props = $props();

	const level = $derived($tezbakeState?.rights?.level ?? 0);
	const timeline = $derived(bakerRightsTimeline($tezbakeState?.rights?.rights ?? [], baker, level));
	const nextBlock = $derived(timeline.upcoming.find((r) => r.blocks > 0)?.level);
	const summary = $derived(summarizeRights(timeline.past));
	const pastBlocks = $derived(summary.blocksDone + summary.blocksMissed);
	const pastAttestations = $derived(summary.attestationsDone + summary.attestationsMissed);
	// the strip shows the most recent past and nearest upcoming levels around "now"
	const strip = $derived([...timeline.past.slice(0, 10).reverse(), ...timeline.upcoming.slice(0, 10)]);
</script>

<Card>
	{#snippet header()}
		<h3 class="sl-card-title">Rights</h3>
		<span class="baker">
			{#if showColor}<span class="color" style:background-color={getBakerColor(baker)}></span>{/if}
			{formatAddress(baker)}
		</span>
	{/snippet}

	<dl class="kv">
		<dt>Next block</dt>
		<dd>
			{#if nextBlock}{nextBlock} <span class="hint">in {nextBlock - level}</span>{:else}<span class="hint"
					>none in window</span
				>{/if}
		</dd>
		<dt>Recent blocks</dt>
		<dd class:text-error={summary.blocksMissed > 0}>{summary.blocksDone}/{pastBlocks}</dd>
		<dt>Recent attestations</dt>
		<dd class:text-error={summary.attestationsMissed > 0}>
			{summary.attestationsDone}/{pastAttestations}
		</dd>
	</dl>

	<RightsStrip rows={strip} {level} />

	{#snippet footer()}
		<Button onclick={() => goto(`/rights?baker=${baker}`)}>Open</Button>
	{/snippet}
</Card>

<style lang="sass">
.baker
	display: flex
	align-items: center
	gap: var(--spacing-f2)
	font-size: 0.85rem
	color: var(--hint-color)

	.color
		width: 10px
		height: 10px
		border-radius: 3px
</style>
