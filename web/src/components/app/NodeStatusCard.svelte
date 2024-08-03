<script lang="ts">
	import { writeToClipboard } from '@src/util/clipboard';
	import Card from '@components/starlight/components/Card.svelte';
	import type { NodeStatus } from '@src/common/types/status';
	import { formatBlockHash, formatTimestampAgoStrict } from '@src/util/format';
	import Separator from './Separator.svelte';
	import { onDestroy } from 'svelte';

	export let node: NodeStatus;
	export let title = `Node`;

	$: blockTimestamp = formatTimestampAgoStrict(node.block?.timestamp ?? 0);

	const interval = setInterval(() => {
		blockTimestamp = formatTimestampAgoStrict(node.block?.timestamp ?? 0);
	}, 500);

	onDestroy(() => clearInterval(interval));
</script>

<Card>
	<div class="node-grid">
		<div class="title">
			<h5>
				{title}
				<div class="connection-status">
					<div
						class="connection-status-sign"
						class:connected={node.connection_status === 'connected'}
					/>
				</div>
			</h5>
			<button class="unstyle-button address" on:click={() => writeToClipboard(node.address ?? '')}
				>{node.address}</button
			>
		</div>
		<Separator />
		{#if node.connection_status === 'connected'}
			<div class="chain-state">
				<div class="level">
					{node.block?.level_info.level}
					<div class="cycle">#{node.block?.level_info.cycle}</div>
				</div>

				<button
					class="unstyle-button hash"
					on:click={() => writeToClipboard(node.block?.hash ?? '')}
					>{formatBlockHash(node.block?.hash ?? '')}</button
				>
				<div class="timestamp">{blockTimestamp}</div>
			</div>
		{:else}
			<div class="disconnected-status">DISCONNECTED</div>
		{/if}
		{#if node.network_info}
			<div class="network-info">
				<div class="connections">
					<div class="value">{node.network_info?.connection_count}</div>
					connections
				</div>
			</div>
		{/if}
	</div>
</Card>

<style lang="sass">
	.node-grid
		display: grid
		grid-template-columns: minmax(100px, 1fr)
		grid-template-rows: auto auto 1fr auto auto 
		height: 100%
		gap: var(--spacing)

		.title
			position: relative

			h5
				font-size: 1.5rem
				font-weight: 500
				margin: 0
			
			.address
				margin: 0
				margin-top: var(--spacing-f2)
				font-size: 0.95rem
		
		.chain-state
			display: grid
			gap: var(--spacing)
			grid-template-rows: auto auto auto 1fr

			.level
				font-size: 1.5rem
				font-weight: 500

				.cycle
					display: inline
					font-size: 1rem

		.bottom-separator
			grid-row: 4
			display: flex
			align-items: flex-end

		.network-info
			grid-row: 5
			.connections
				.value
					font-size: 1.25rem
					font-weight: 500
					display: inline-block


.disconnected-status
	display: flex
	align-items: center
	justify-content: center
	font-size: 1.25rem
	height: 100%

.connection-status
	.connection-status-sign
		position: absolute
		display: inline-block
		width: 12px
		height: 12px
		border-radius: 50%
		margin-right: var(--spacing-f2)
		background-color: var(--error-color)
		top: 0
		right: 0

		&.connected
			background-color: var(--success-color) !important

</style>
