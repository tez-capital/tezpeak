<script lang="ts">
	import { Card } from '@alis-is/starlight';
	import { writeToClipboard } from '@src/util/clipboard';
	import type { NodeStatus } from '@src/common/types/status';
	import { formatBlockHash, formatTimestampAgoStrict } from '@src/util/format';
	import { now } from '@app/state/clock';

	interface Props {
		node: NodeStatus;
		title?: string;
	}
	let { node, title = 'Node' }: Props = $props();

	const connected = $derived(node.connection_status === 'connected');
</script>

<Card>
	{#snippet header()}
		<h3 class="sl-card-title">{title}</h3>
		<span class="state">
			{#if !node.is_essential}<span class="hint">optional</span>{/if}
			<span class="status-dot" class:ok={connected} class:error={!connected} title={node.connection_status}></span>
		</span>
	{/snippet}

	<button class="unstyle-button address" title="Copy address" onclick={() => writeToClipboard(node.address ?? '')}
		>{node.address}</button
	>
	{#if connected && node.block}
		<div class="metric">{node.block.level_info.level} <small>cycle {node.block.level_info.cycle}</small></div>
		<dl class="kv">
			<dt>Block</dt>
			<dd>
				<button class="unstyle-button" title="Copy block hash" onclick={() => writeToClipboard(node.block?.hash ?? '')}
					>{formatBlockHash(node.block.hash)}</button
				>
			</dd>
			<dt>Received</dt>
			<dd>{#key $now}{formatTimestampAgoStrict(node.block.timestamp)}{/key}</dd>
			{#if node.network_info}
				<dt>Connections</dt>
				<dd>{node.network_info.connection_count}</dd>
			{/if}
		</dl>
	{:else}
		<div class="disconnected" class:text-error={node.is_essential} class:text-warn={!node.is_essential}>
			{node.connection_status}
		</div>
	{/if}
</Card>

<style lang="sass">
.state
	display: flex
	align-items: center
	gap: var(--spacing-f2)

.disconnected
	align-self: center
	padding: var(--spacing) 0
	text-align: center
	text-transform: uppercase
	letter-spacing: 0.1rem
</style>
