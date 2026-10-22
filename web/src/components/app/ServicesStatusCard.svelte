<script lang="ts">
	import { Badge, Card } from '@alis-is/starlight';
	import type { ServicesStatus } from '@src/common/types/status';
	import { formatTimestamp, formatTimestampAgo } from '@src/util/format';
	import { now } from '@app/state/clock';

	interface Props {
		title?: string;
		services?: ServicesStatus;
	}
	let { title = 'Services', services = { timestamp: 0, applications: {} } }: Props = $props();

	const applications = $derived(
		Object.entries(services.applications ?? {}).map(([id, appServices]) => [id, Object.entries(appServices)] as const)
	);
	const all = $derived(applications.flatMap(([, s]) => s));
	const runningCount = $derived(all.filter(([, s]) => s.status === 'running').length);
</script>

<Card>
	{#snippet header()}
		<h3 class="sl-card-title">{title}</h3>
		<Badge variant={runningCount === all.length ? 'success' : 'error'}>{runningCount}/{all.length} running</Badge>
	{/snippet}

	<div class="services">
		{#each applications as [app, appServices] (app)}
			<div class="app">{app}</div>
			{#each appServices as [id, info] (id)}
				{@const running = info.status === 'running'}
				<span class="status-dot" class:ok={running} class:error={!running}></span>
				<span class="id">{id}</span>
				<span class="state" class:text-error={!running}>
					{#key $now}{running ? formatTimestampAgo(info.started) : info.status}{/key}
				</span>
			{:else}
				<span class="hint none">no services</span>
			{/each}
		{/each}
	</div>

	{#snippet footer()}
		<span class="hint">updated {formatTimestamp(services.timestamp)}</span>
	{/snippet}
</Card>

<style lang="sass">
.services
	display: grid
	grid-template-columns: auto minmax(0, 1fr) auto
	align-items: center
	column-gap: var(--spacing-f2)
	row-gap: 6px

.app
	grid-column: 1 / -1
	margin-top: var(--spacing-f2)
	font-size: 0.8rem
	letter-spacing: 0.08rem
	text-transform: uppercase
	color: var(--hint-color)

	&:first-child
		margin-top: 0

.id
	text-transform: capitalize
	overflow: hidden
	text-overflow: ellipsis
	white-space: nowrap

.state
	text-align: right
	white-space: nowrap
	font-variant-numeric: tabular-nums

.none
	grid-column: 2 / -1
</style>
