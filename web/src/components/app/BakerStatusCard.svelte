<script lang="ts">
	import { Badge, Card } from '@alis-is/starlight';
	import { writeToClipboard } from '@src/util/clipboard';
	import { formatBalance } from '@src/util/format';
	import { getBakerColor } from '@src/util/baker';
	import type { BakerStatus } from '@src/common/types/status';

	interface Props {
		baker: string;
		status: BakerStatus;
		showColor?: boolean;
	}
	let { baker, status, showColor = false }: Props = $props();
</script>

<Card>
	{#snippet header()}
		<h3 class="sl-card-title">Baker</h3>
		{#if status.deactivated}
			<Badge variant="error">deactivated</Badge>
		{:else if showColor}
			<span class="color" style:background-color={getBakerColor(baker)} title="baker colour"></span>
		{/if}
	{/snippet}

	<button class="unstyle-button address" title="Copy address" onclick={() => writeToClipboard(baker)}
		>{baker}</button
	>
	<dl class="kv">
		<dt>Balance</dt>
		<dd>{formatBalance(status.balance)}</dd>
		<dt>Staked</dt>
		<dd>{formatBalance(status.staked_balance)}</dd>
		<dt>External staked</dt>
		<dd>{formatBalance(status.external_staked_balance)}</dd>
		<dt>External delegated</dt>
		<dd>{formatBalance(status.external_delegated_balance)}</dd>
	</dl>

	{#snippet footer()}
		<div class="metric">{status.delegators_count} <small>delegators</small></div>
	{/snippet}
</Card>

<style lang="sass">
.color
	width: 10px
	height: 10px
	border-radius: 3px
</style>
