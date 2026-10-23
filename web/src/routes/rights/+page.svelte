<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { Button, Card, Tabs, Toggle } from '@alis-is/starlight';
	import HomeIcon from '@alis-is/starlight/la/icons/home-solid.svelte';
	import { bakers, state as tezbakeState } from '@app/state/tezbake';
	import { now } from '@app/state/clock';
	import RightsLanes from '@components/app/RightsLanes.svelte';
	import {
		bakerRights,
		bakerRightsAt,
		getBakerColor,
		summarizeRights,
		worstStatus,
		type BakerLevelRights,
		type RightStatus
	} from '@src/util/baker';
	import { formatAddress, formatBlocksEta, formatShortAddress } from '@src/util/format';

	const ALL = 'all';

	const bakerIds = $derived($bakers.map(([id]) => id));
	const requested = $derived(page.url.searchParams.get('baker') ?? '');
	// combined view by default when there is more than one baker
	const selected = $derived(
		bakerIds.includes(requested) ? requested : bakerIds.length > 1 ? ALL : (bakerIds[0] ?? '')
	);
	const shown = $derived(selected === ALL ? bakerIds : [selected]);
	const combined = $derived(shown.length > 1);

	const level = $derived($tezbakeState?.rights?.level ?? 0);
	const rights = $derived($tezbakeState?.rights?.rights ?? []);
	// newest first, upcoming levels above "now"; consecutive routine levels (no block, nothing missed, same status
	// per baker) fold into one row with summed slots, so blocks and misses stand out
	let everyLevel = $state(false);
	type Row = {
		from: number;
		to: number;
		count: number;
		cells: Array<BakerLevelRights | undefined>;
	};
	const rows = $derived.by(() => {
		const result = [] as Array<Row & { key: string }>;
		for (const r of [...rights].sort((a, b) => b.level - a.level)) {
			const cells = shown.map((b) => bakerRightsAt(r, b, level));
			const routine =
				!everyLevel && cells.every((c) => !c || (!c.blocks && worstStatus(c) !== 'missed'));
			const key = routine
				? `${r.level > level}:${cells.map((c) => (c ? worstStatus(c) : '-'))}`
				: '';
			const last = result.at(-1);
			if (key && last?.key === key) {
				last.to = r.level;
				last.count++;
				last.cells = last.cells.map(
					(c, i) => c && { ...c, attestations: c.attestations + (cells[i]?.attestations ?? 0) }
				);
			} else {
				result.push({ from: r.level, to: r.level, count: 1, cells, key });
			}
		}
		return result;
	});
	const nowIndex = $derived(rows.findIndex((r) => r.from <= level));
	const span = (r: Row) =>
		r.to > level
			? `in ${r.to - level}–${r.from - level}`
			: `${r.from === level ? 'now' : level - r.from}–${level - r.to} ago`;
	// 9123549–67, like page ranges
	function levelRange(r: Row) {
		const [to, from] = [String(r.to), String(r.from)];
		let same = 0;
		while (same < from.length - 2 && to[same] === from[same]) same++;
		return `${to}–${from.slice(same)}`;
	}
	const pastCount = $derived(rights.filter((r) => r.level <= level).length);

	const perBaker = $derived(
		shown.map((baker) => {
			const rows = bakerRights(rights, baker, level);
			return { baker, rows, summary: summarizeRights(rows) };
		})
	);
	const byLevel = (a: { level: number }, b: { level: number }) => a.level - b.level;
	const upcomingBlocks = $derived(
		perBaker
			.flatMap(({ baker, rows }) =>
				rows.filter((r) => r.blockStatus === 'upcoming').map((r) => ({ ...r, baker }))
			)
			.sort(byLevel)
	);
	const missed = $derived(
		perBaker
			.flatMap(({ baker, rows }) =>
				rows
					.filter((r) => r.blockStatus === 'missed' || r.attestationStatus === 'missed')
					.map((r) => ({ ...r, baker }))
			)
			.sort((a, b) => b.level - a.level)
	);

	function select(baker: string) {
		goto(`?baker=${baker}`, { replaceState: true, noScroll: true, keepFocus: true });
	}

	const WORD: Record<RightStatus, string> = {
		upcoming: 'scheduled',
		pending: 'checking',
		done: 'ok',
		missed: 'missed'
	};
	const relative = (l: number) =>
		l > level ? `in ${l - level}` : l === level ? 'now' : `${level - l} ago`;
</script>

{#snippet bakerTag(baker: string)}
	<span class="baker-tag" title={baker}>
		<span class="swatch" style:background-color={getBakerColor(baker)}></span>{formatShortAddress(
			baker
		)}
	</span>
{/snippet}

{#snippet cell(r: BakerLevelRights | undefined)}
	{#if r}
		{#if r.blockStatus}
			<span class="block status-{r.blockStatus}"
				>block {r.blockStatus === 'done' ? 'baked' : WORD[r.blockStatus]}</span
			>
		{/if}
		{#if r.attestationStatus}
			<span class="status-{r.attestationStatus}">
				{r.attestations} <span class="hint slots">slots</span>
				{r.attestationStatus === 'done' ? 'attested' : WORD[r.attestationStatus]}
			</span>
		{/if}
	{:else}
		<span class="hint">—</span>
	{/if}
{/snippet}

<div class="rights-page">
	<nav class="page-navigation">
		<Button onclick={() => goto('/')}>
			<div class="navigation-btn-content"><HomeIcon /> HOME</div>
		</Button>
	</nav>

	{#if bakerIds.length === 0}
		<Card><p class="hint">No bakers are tracked.</p></Card>
	{:else if !rights.length}
		<Card
			><p class="hint">
				No rights data yet. Rights are tracked once the node reports a new block.
			</p></Card
		>
	{:else}
		{#if bakerIds.length > 1}
			<Tabs
				tabs={[
					{ label: 'All bakers', value: ALL },
					...bakerIds.map((id) => ({ label: formatAddress(id), value: id }))
				]}
				value={selected}
				onchange={select}
			/>
		{/if}

		<Card>
			{#snippet header()}
				<h3 class="sl-card-title">Rights window</h3>
				{#if combined}
					<span class="hint">{pastCount} past · {rights.length - pastCount} upcoming levels</span>
				{:else}
					{@const s = perBaker[0].summary}
					<span class="hint">
						<span class:text-error={s.attestationsMissed > 0}
							>{s.attestationsDone}/{s.attestationsDone + s.attestationsMissed} attested</span
						>
						{#if s.blocksDone + s.blocksMissed}
							· <span class:text-error={s.blocksMissed > 0}
								>{s.blocksDone}/{s.blocksDone + s.blocksMissed} baked</span
							>{/if}
					</span>
				{/if}
			{/snippet}
			<RightsLanes {rights} {level} bakers={shown} labels={combined} />
			<div class="legend hint">
				<span><i class="dot done"></i>attested</span>
				<span><i class="dot missed"></i>missed</span>
				<span><i class="dot pending"></i>checking</span>
				<span><i class="dot upcoming"></i>scheduled</span>
				<span><i class="dot done tall"></i>block</span>
			</div>
		</Card>

		<div class="columns">
			{#if combined}
				<Card>
					{#snippet header()}
						<h3 class="sl-card-title">Performance</h3>
						<span class="hint">last {pastCount} levels</span>
					{/snippet}
					<div class="table perf">
						<div class="head">
							<div>Baker</div>
							<div>Attested</div>
							<div>Baked</div>
						</div>
						{#each perBaker as { baker, summary } (baker)}
							{@const attestations = summary.attestationsDone + summary.attestationsMissed}
							{@const blocks = summary.blocksDone + summary.blocksMissed}
							<div>{@render bakerTag(baker)}</div>
							<div class:text-error={summary.attestationsMissed > 0}>
								{summary.attestationsDone}/{attestations}
								{#if attestations}<span class="hint"
										>{Math.floor((summary.attestationsDone / attestations) * 100)}%</span
									>{/if}
							</div>
							<div class:text-error={summary.blocksMissed > 0}>
								{#if blocks}{summary.blocksDone}/{blocks}{:else}<span class="hint">none due</span
									>{/if}
							</div>
						{/each}
					</div>
				</Card>
			{/if}

			<Card>
				{#snippet header()}
					<h3 class="sl-card-title">Upcoming blocks</h3>
					<span class="hint">next {rights.length - pastCount} levels</span>
				{/snippet}
				{#each upcomingBlocks as r (`${r.baker}-${r.level}`)}
					{@const eta = formatBlocksEta(r.level - level, $now)}
					<div class="event">
						<span class="level">{r.level}</span>
						{#if combined}{@render bakerTag(r.baker)}{/if}
						<span class="when">{eta.in} <span class="hint">· {eta.at}</span></span>
					</div>
				{:else}
					<span class="hint">No block rights in the window.</span>
				{/each}
			</Card>

			<Card>
				{#snippet header()}
					<h3 class="sl-card-title">Missed</h3>
					<span class="hint">last {pastCount} levels</span>
				{/snippet}
				{#each missed as r (`${r.baker}-${r.level}`)}
					<div class="event">
						<span class="level">{r.level}</span>
						{#if combined}{@render bakerTag(r.baker)}{/if}
						<span class="text-error"
							>{r.blockStatus === 'missed'
								? 'block'
								: `attestation · ${r.attestations} slots`}</span
						>
						<span class="when hint">{relative(r.level)}</span>
					</div>
				{:else}
					<span class="text-ok">Nothing missed.</span>
				{/each}
			</Card>
		</div>

		<Card>
			{#snippet header()}
				<h3 class="sl-card-title">Levels</h3>
				<label class="hint every">every level <Toggle bind:checked={everyLevel} /></label>
			{/snippet}
			<div class="table levels" style:--bakers={shown.length} role="table">
				<div class="head" role="row">
					<div>Level</div>
					{#each shown as baker (baker)}<div>
							{#if combined}{@render bakerTag(baker)}{:else}Rights{/if}
						</div>{/each}
				</div>
				{#each rows as r, i (r.from)}
					{#if i === nowIndex}<div class="now" role="row"><span>now</span></div>{/if}
					<div
						class="row"
						class:has-block={r.cells.some((c) => c && c.blocks > 0)}
						class:future={r.to > level}
						role="row"
					>
						<div class="level">
							{#if r.from === r.to}{r.from}
								<div class="hint">{relative(r.from)}</div>
							{:else}<span title="{r.count} levels">{levelRange(r)}</span>
								<div class="hint">{span(r)} · {r.count} levels</div>
							{/if}
						</div>
						{#each r.cells as c, i (shown[i])}<div class="cell">{@render cell(c)}</div>{/each}
					</div>
				{/each}
				<!-- every row upcoming: now is below them -->
				{#if nowIndex === -1}<div class="now" role="row"><span>now</span></div>{/if}
			</div>
		</Card>
	{/if}
</div>

<style lang="sass">
.rights-page
	display: grid
	grid-template-columns: minmax(0, 1fr)
	gap: var(--spacing)

.swatch
	flex: none
	width: 10px
	height: 10px
	border-radius: 3px

.baker-tag
	display: inline-flex
	align-items: center
	gap: var(--spacing-f2)
	white-space: nowrap
	text-transform: none
	letter-spacing: normal

.every
	display: flex
	align-items: center
	gap: var(--spacing-f2)
	white-space: nowrap

.legend
	display: flex
	flex-wrap: wrap
	gap: var(--spacing-f2) var(--spacing)

	span
		display: inline-flex
		align-items: center
		gap: 6px

.dot
	width: 8px
	height: 8px
	border-radius: 50%

	&.tall
		width: 6px
		height: 14px
		border-radius: 2px

	&.done
		background-color: var(--success-color)
	&.missed
		background-color: var(--error-color)
	&.pending
		background-color: var(--hint-color)
		opacity: 0.6
	&.upcoming
		box-shadow: inset 0 0 0 1.5px var(--accent)

.columns
	display: grid
	grid-template-columns: repeat(auto-fit, minmax(min(100%, 300px), 1fr))
	align-items: start
	gap: var(--spacing)

.event
	display: flex
	align-items: baseline
	flex-wrap: wrap
	gap: 0 var(--spacing)
	font-variant-numeric: tabular-nums

	.when
		margin-left: auto
		white-space: nowrap

.level
	font-variant-numeric: tabular-nums

.table
	display: grid
	font-variant-numeric: tabular-nums

	> div, .row > div
		padding: 6px var(--spacing) 6px 0
		border-bottom: 1px solid var(--frost-border-color)

	.head, .row
		display: contents

	.head > div
		font-size: 0.8rem
		letter-spacing: 0.08rem
		text-transform: uppercase
		color: var(--hint-color)

.perf
	grid-template-columns: auto auto auto

	.hint
		margin-left: var(--spacing-f2)

.levels
	grid-template-columns: auto repeat(var(--bakers), minmax(0, 1fr))

	.row.future > div
		color: var(--hint-color)

	// block rights are rare (a few per cycle), so their rows stand out
	.row.has-block > div
		background-color: color-mix(in srgb, var(--accent) 12%, transparent)

		&.level
			font-weight: 600
			box-shadow: inset 3px 0 0 var(--accent)
			padding-left: var(--spacing-f2)

	.cell
		display: flex
		flex-wrap: wrap
		gap: 0 var(--spacing)

	.block
		font-weight: 600

	.now
		grid-column: 1 / -1
		display: flex
		align-items: center
		gap: var(--spacing-f2)
		border-bottom: none
		color: var(--accent)
		font-size: 0.85rem

		&::before, &::after
			content: ''
			flex: 1
			height: 2px
			border-radius: 1px
			background-color: var(--accent)
			opacity: 0.6

.status-done
	color: var(--success-color)
.status-missed
	color: var(--error-color)
.status-upcoming
	color: var(--accent)

@media (max-width: 600px)
	.levels .slots
		display: none
</style>
