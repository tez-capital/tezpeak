<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { Badge, Button, Card, Tabs } from '@alis-is/starlight';
	import HomeIcon from '@alis-is/starlight/la/icons/home-solid.svelte';
	import { bakers, state as tezbakeState } from '@app/state/tezbake';
	import RightsStrip from '@components/app/RightsStrip.svelte';
	import {
		bakerRightsTimeline,
		getBakerColor,
		mergeByLevel,
		summarizeRights,
		type BakerLevelRights,
		type RightStatus
	} from '@src/util/baker';
	import { formatAddress } from '@src/util/format';
	import { writeToClipboard } from '@src/util/clipboard';

	const ALL = 'all';
	type Row = BakerLevelRights & { baker: string };

	const bakerIds = $derived($bakers.map(([id]) => id));
	const requested = $derived(page.url.searchParams.get('baker') ?? '');
	// combined view by default when there is more than one baker
	const selected = $derived(
		bakerIds.includes(requested) || (requested === ALL && bakerIds.length > 1)
			? requested
			: bakerIds.length > 1
				? ALL
				: (bakerIds[0] ?? '')
	);
	const combined = $derived(selected === ALL);
	const shown = $derived(combined ? bakerIds : [selected]);

	const level = $derived($tezbakeState?.rights?.level ?? 0);
	const timelines = $derived(
		shown.map((baker) => {
			const t = bakerRightsTimeline($tezbakeState?.rights?.rights ?? [], baker, level);
			const tag = (r: BakerLevelRights): Row => ({ ...r, baker });
			return { upcoming: t.upcoming.map(tag), past: t.past.map(tag) };
		})
	);
	const byLevelThenBaker = (dir: 1 | -1) => (a: Row, b: Row) =>
		(a.level - b.level) * dir || bakerIds.indexOf(a.baker) - bakerIds.indexOf(b.baker);
	const upcoming = $derived(timelines.flatMap((t) => t.upcoming).sort(byLevelThenBaker(1)));
	const past = $derived(timelines.flatMap((t) => t.past).sort(byLevelThenBaker(-1)));

	const summary = $derived(summarizeRights(past));
	const upcomingBlocks = $derived(upcoming.filter((r) => r.blocks > 0).length);
	const upcomingSlots = $derived(upcoming.reduce((sum, r) => sum + r.attestations, 0));
	const strip = $derived(mergeByLevel([...[...past].reverse(), ...upcoming]));

	function select(baker: string) {
		goto(`?baker=${baker}`, { replaceState: true, noScroll: true, keepFocus: true });
	}

	const BADGE: Record<RightStatus, 'accent' | 'default' | 'success' | 'error'> = {
		upcoming: 'accent',
		pending: 'default',
		done: 'success',
		missed: 'error'
	};
	const LABEL: Record<RightStatus, string> = {
		upcoming: 'scheduled',
		pending: 'checking',
		done: 'ok',
		missed: 'missed'
	};
	const BLOCK_LABEL: Record<RightStatus, string> = { ...LABEL, done: 'baked' };
	const ATTESTATION_LABEL: Record<RightStatus, string> = { ...LABEL, done: 'attested' };
</script>

{#snippet bakerTag(baker: string)}
	<span class="baker-tag" title={baker}>
		<span class="color" style:background-color={getBakerColor(baker)}></span>{formatAddress(baker)}
	</span>
{/snippet}

{#snippet row(r: Row, sameLevel: boolean)}
	<!-- further bakers on the same level leave level/when blank so a level reads as one group -->
	<div class="level">{sameLevel ? '' : r.level}</div>
	<div class="when hint">
		{#if !sameLevel}{r.level > level ? `in ${r.level - level}` : r.level === level ? 'now' : `${level - r.level} ago`}{/if}
	</div>
	{#if combined}<div class="who">{@render bakerTag(r.baker)}</div>{/if}
	<div class="block">
		{#if r.blockStatus}<Badge variant={BADGE[r.blockStatus]}>{BLOCK_LABEL[r.blockStatus]}</Badge>{/if}
	</div>
	<div class="attestation">
		{#if r.attestationStatus}
			<span class="slots">{r.attestations}</span>
			<Badge variant={BADGE[r.attestationStatus]} outline>{ATTESTATION_LABEL[r.attestationStatus]}</Badge>
		{/if}
	</div>
{/snippet}

<div class="rights-page">
	<nav class="page-navigation">
		<Button onclick={() => goto('/')}>
			<div class="navigation-btn-content"><HomeIcon /> HOME</div>
		</Button>
	</nav>

	{#if bakerIds.length === 0}
		<Card><p class="hint">No bakers are tracked.</p></Card>
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
				<h3 class="sl-card-title">Baking rights</h3>
				{#if combined}
					<span class="hint">{bakerIds.length} bakers combined</span>
				{:else}
					<button class="unstyle-button baker" title="Copy {selected}" onclick={() => writeToClipboard(selected)}>
						<span class="color" style:background-color={getBakerColor(selected)}></span>
						{formatAddress(selected)}
					</button>
				{/if}
			{/snippet}

			<div class="metrics">
				<div><div class="metric">{level}</div><div class="hint">current level</div></div>
				<div><div class="metric">{upcomingBlocks}</div><div class="hint">upcoming blocks</div></div>
				<div><div class="metric">{upcomingSlots}</div><div class="hint">upcoming attestation slots</div></div>
				<div>
					<div class="metric" class:text-error={summary.blocksMissed > 0}>
						{summary.blocksDone}<small>/{summary.blocksDone + summary.blocksMissed}</small>
					</div>
					<div class="hint">blocks baked</div>
				</div>
				<div>
					<div class="metric" class:text-error={summary.attestationsMissed > 0}>
						{summary.attestationsDone}<small>/{summary.attestationsDone + summary.attestationsMissed}</small>
					</div>
					<div class="hint">levels attested</div>
				</div>
			</div>
			<RightsStrip rows={strip} {level} />
		</Card>

		<Card>
			{#snippet header()}
				<h3 class="sl-card-title">Timeline</h3>
				<span class="hint">levels with rights in the tracked window</span>
			{/snippet}

			<div class="timeline" class:combined role="table">
				<div class="head" role="row">
					<div>Level</div>
					<div class="when">When</div>
					{#if combined}<div class="who">Baker</div>{/if}
					<div>Block</div>
					<div>Attestation slots</div>
				</div>

				<div class="group" role="row">Upcoming · {upcoming.length}</div>
				{#each upcoming as r, i (`${r.baker}-${r.level}`)}{@render row(r, i > 0 && upcoming[i - 1].level === r.level)}{:else}
					<div class="empty hint">no upcoming rights in the window</div>
				{/each}

				<div class="now-marker" role="row"><span>now · level {level}</span></div>

				<div class="group" role="row">Past · {past.length}</div>
				{#each past as r, i (`${r.baker}-${r.level}`)}{@render row(r, i > 0 && past[i - 1].level === r.level)}{:else}
					<div class="empty hint">no past rights in the window</div>
				{/each}
			</div>
		</Card>
	{/if}
</div>

<style lang="sass">
.rights-page
	display: grid
	grid-template-columns: minmax(0, 1fr)
	gap: var(--spacing)

.baker
	display: flex
	align-items: center
	gap: var(--spacing-f2)
	font-size: 0.85rem
	color: var(--hint-color)
	overflow-wrap: anywhere
	text-align: right

.color
	flex: none
	width: 10px
	height: 10px
	border-radius: 3px

.baker-tag
	display: inline-flex
	align-items: center
	gap: var(--spacing-f2)
	font-size: 0.85rem
	white-space: nowrap

.metrics
	display: grid
	grid-template-columns: repeat(auto-fit, minmax(140px, 1fr))
	gap: var(--spacing)
	padding-bottom: var(--spacing-f2)

.timeline
	display: grid
	grid-template-columns: minmax(90px, auto) minmax(80px, auto) minmax(110px, 1fr) minmax(170px, 2fr)
	align-items: stretch
	font-variant-numeric: tabular-nums

	&.combined
		grid-template-columns: minmax(90px, auto) minmax(80px, auto) minmax(170px, auto) minmax(110px, 1fr) minmax(170px, 2fr)

	// no column gap: cell padding keeps the row borders continuous
	> div
		display: flex
		align-items: center
		min-height: 24px
		padding: 6px var(--spacing) 6px 0
		border-bottom: 1px solid var(--frost-border-color)

	.head
		display: contents

		> div
			padding: 0 0 var(--spacing-f2)
			font-size: 0.8rem
			letter-spacing: 0.08rem
			text-transform: uppercase
			color: var(--hint-color)
			border-bottom: 1px solid var(--frost-border-color)

	.group, .empty, .now-marker
		grid-column: 1 / -1

	.group
		padding-top: var(--spacing)
		font-size: 0.8rem
		letter-spacing: 0.08rem
		text-transform: uppercase
		color: var(--hint-color)

	.attestation
		display: flex
		align-items: center
		gap: var(--spacing-f2)

		.slots
			min-width: 2ch
			text-align: right

	.now-marker
		display: flex
		align-items: center
		gap: var(--spacing-f2)
		padding: var(--spacing-f2) 0
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

@media (max-width: 600px)
	.timeline, .timeline.combined
		grid-template-columns: auto minmax(0, 1fr) minmax(0, 1.4fr)

		.when
			display: none

	// on phones the baker shows as a coloured mark next to the level
	.timeline.combined
		grid-template-columns: auto auto minmax(0, 1fr) minmax(0, 1.4fr)

		.baker-tag
			font-size: 0
			gap: 0

		.head .who
			font-size: 0
</style>
