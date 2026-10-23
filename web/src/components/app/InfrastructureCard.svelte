<script lang="ts">
	import { Card } from '@alis-is/starlight';
	import type { LedgerWalletStatus, NodeStatus, ServicesStatus } from '@src/common/types/status';
	import { formatTimestamp, formatTimestampAgoStrict } from '@src/util/format';
	import { now } from '@app/state/clock';

	interface Props {
		nodes: Array<[string, NodeStatus]>;
		services?: ServicesStatus;
		wallets: Array<[string, LedgerWalletStatus]>;
	}
	let { nodes, services, wallets }: Props = $props();

	const allServices = $derived(
		Object.values(services?.applications ?? {}).flatMap((app) => Object.entries(app))
	);
	const running = $derived(allServices.filter(([, s]) => s.status === 'running').length);
</script>

<Card>
	{#snippet header()}
		<h3 class="sl-card-title">Infrastructure</h3>
	{/snippet}

	<div class="group">Nodes</div>
	{#each nodes as [id, node] (id)}
		{@const connected = node.connection_status === 'connected'}
		<div class="row">
			<span
				class="status-dot"
				class:ok={connected}
				class:error={!connected && node.is_essential}
				class:warn={!connected && !node.is_essential}
			></span>
			<div class="name" title={node.address}>{id}</div>
			<div class="value">
				{#if connected && node.block}
					<span title={node.block.hash}>{node.block.level_info.level}</span>
					<div class="hint">
						{#key $now}{formatTimestampAgoStrict(
								node.block.timestamp
							)}{/key}{#if node.network_info}{` · ${node.network_info.connection_count} peers`}{/if}
					</div>
				{:else}
					<span class:text-error={node.is_essential} class:text-warn={!node.is_essential}
						>{node.connection_status}</span
					>
					{#if !node.is_essential}<div class="hint">optional</div>{/if}
				{/if}
			</div>
		</div>
	{:else}
		<span class="hint">no nodes configured</span>
	{/each}

	{#if allServices.length}
		<div class="group">
			Services
			<span class:text-error={running < allServices.length}
				>{running}/{allServices.length} running</span
			>
		</div>
		<div class="chips">
			{#each allServices as [id, s] (id)}
				{@const ok = s.status === 'running'}
				<span
					class="chip"
					class:failed={!ok}
					title={ok ? `running since ${formatTimestamp(s.started)}` : s.status}
				>
					<span class="status-dot" class:ok class:error={!ok}></span>{id}{#if !ok}<span
							>· {s.status}</span
						>{/if}
				</span>
			{/each}
		</div>
	{/if}

	{#if wallets.length}
		<div class="group">Ledger</div>
		{#each wallets as [id, w] (id)}
			{@const ready = w.ledger_status === 'connected' && w.authorized}
			<div class="row">
				<span class="status-dot" class:ok={ready} class:error={!ready}></span>
				<div class="name">
					{id}
					<div class="hint">
						{w.ledger || 'no device'}{#if w.app_version}{` · app ${w.app_version}`}{/if}
					</div>
				</div>
				<div class="value" class:text-error={!ready}>
					{w.ledger_status !== 'connected'
						? w.ledger_status || 'disconnected'
						: w.authorized
							? 'ready'
							: 'not authorized'}
				</div>
			</div>
		{/each}
	{/if}
</Card>

<style lang="sass">
.group
	display: flex
	justify-content: space-between
	margin-top: var(--spacing-f2)
	font-size: 0.8rem
	letter-spacing: 0.08rem
	text-transform: uppercase
	color: var(--hint-color)

	&:first-child
		margin-top: 0

	span
		letter-spacing: normal
		text-transform: none

.row
	display: grid
	grid-template-columns: auto minmax(0, 1fr) auto
	align-items: baseline
	column-gap: var(--spacing-f2)

.name
	min-width: 0


.value
	text-align: right
	white-space: nowrap
	font-variant-numeric: tabular-nums

.chips
	display: flex
	flex-wrap: wrap
	gap: 6px

.chip
	display: inline-flex
	align-items: center
	gap: 6px
	padding: 2px 10px
	border-radius: 999px
	border: 1px solid var(--frost-border-color)
	font-size: 0.85rem

	&.failed
		border-color: var(--error-color)
		color: var(--error-color)
</style>
