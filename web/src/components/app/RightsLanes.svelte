<script lang="ts">
	import type { BlockRights } from '@src/common/types/status';
	import { bakerRightsAt, describeRights, getBakerColor, worstStatus } from '@src/util/baker';
	import { formatShortAddress } from '@src/util/format';

	interface Props {
		rights: Array<BlockRights>;
		level: number;
		bakers: Array<string>;
		// baker label in front of each lane
		labels?: boolean;
	}
	let { rights, level, bakers, labels = false }: Props = $props();

	// narrow screens show only the levels nearest to "now", ~11px each, split evenly around it when possible
	let width = $state(0);
	const fit = $derived(
		width ? Math.max(12, Math.floor((width - (labels ? 120 : 0)) / 11)) : Infinity
	);
	const sorted = $derived.by(() => {
		const all = [...rights].sort((a, b) => a.level - b.level);
		const past = all.filter((r) => r.level <= level);
		const future = all.slice(past.length);
		const futureShown = Math.min(future.length, Math.max(Math.floor(fit / 2), fit - past.length));
		const pastShown = Math.min(past.length, fit - futureShown);
		return [...past.slice(past.length - pastShown), ...future.slice(0, futureShown)];
	});
	const levels = $derived(sorted.map((r) => r.level));
	const pastCount = $derived(levels.filter((l) => l <= level).length);
	const relative = (l: number) =>
		l > level ? `in ${l - level}` : l === level ? 'now' : `${level - l} ago`;
</script>

<!-- one lane per baker, one column per level; the "now" line separates checked past levels from upcoming ones -->
<div
	class="lanes"
	class:labels
	bind:clientWidth={width}
	role="img"
	aria-label="baking rights from level {levels[0]} to {levels.at(-1)}"
>
	{#each bakers as baker (baker)}
		<div class="lane">
			{#if labels}
				<span class="label" title={baker}>
					<span class="swatch" style:background-color={getBakerColor(baker)}
					></span>{formatShortAddress(baker)}
				</span>
			{/if}
			<div class="cells">
				{#each sorted as r, i (r.level)}
					{@const cell = bakerRightsAt(r, baker, level)}
					{#if i === pastCount}<span class="now"></span>{/if}
					<span
						class="cell"
						title="level {r.level} · {relative(r.level)}{cell
							? ` — ${describeRights(cell)}`
							: ' — no rights'}"
					>
						<i
							class="glyph {cell ? worstStatus(cell) : 'none'}"
							class:block={cell && cell.blocks > 0}
						></i>
					</span>
				{/each}
				{#if pastCount === sorted.length}<span class="now"></span>{/if}
			</div>
		</div>
	{/each}
	{#if levels.length}
		<div class="lane axis hint">
			{#if labels}<span class="label"></span>{/if}
			<div class="cells">
				<span style:flex={pastCount}>{levels[0]}</span>
				<span class="now-label">now · {level}</span>
				<span style:flex={sorted.length - pastCount} class="end">{levels.at(-1)}</span>
			</div>
		</div>
	{:else}
		<span class="hint">no rights in the tracked window</span>
	{/if}
</div>

<style lang="sass">
.lanes
	display: grid
	gap: 6px
	font-variant-numeric: tabular-nums

.lane
	display: flex
	align-items: center
	gap: var(--spacing)

.label
	flex: none
	width: 112px
	display: flex
	align-items: center
	gap: var(--spacing-f2)
	font-size: 0.85rem
	white-space: nowrap
	overflow: hidden

.swatch
	flex: none
	width: 10px
	height: 10px
	border-radius: 3px

.cells
	flex: 1
	min-width: 0
	display: flex
	align-items: center
	gap: 2px

.cell
	flex: 1 1 0
	min-width: 0
	height: 22px
	display: flex
	align-items: center
	justify-content: center

.glyph
	width: min(100%, 8px)
	aspect-ratio: 1
	border-radius: 50%
	background-color: var(--success-color)

	// block rights are rare, so they get the tall mark
	&.block
		width: min(100%, 10px)
		height: 20px
		aspect-ratio: auto
		border-radius: 3px

	&.missed
		background-color: var(--error-color)
	&.pending
		background-color: var(--hint-color)
		opacity: 0.6
	&.upcoming
		background-color: transparent
		box-shadow: inset 0 0 0 1.5px var(--accent)
	&.none
		width: 3px
		background-color: var(--frost-border-color)

.now
	flex: none
	align-self: stretch
	width: 2px
	margin: 0 3px
	border-radius: 1px
	background-color: var(--accent)

.axis
	font-size: 0.75rem

	.cells
		gap: 0

	span
		white-space: nowrap

	.end
		text-align: right

	// as wide as the marker above, text centred on it and overflowing both sides
	.now-label
		flex: none
		width: 8px
		display: flex
		justify-content: center
		color: var(--accent)

@media (max-width: 600px)
	.lane
		gap: var(--spacing-f2)

	.label
		width: 92px
		font-size: 0.75rem
</style>
