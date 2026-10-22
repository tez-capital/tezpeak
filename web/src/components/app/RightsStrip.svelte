<script lang="ts">
	import type { BakerLevelRights, RightStatus } from '@src/util/baker';

	interface Props {
		// ascending by level
		rows: Array<BakerLevelRights>;
		level: number;
	}
	let { rows, level }: Props = $props();

	// one colour per level: any miss wins, then pending, then upcoming
	function combined(r: BakerLevelRights): RightStatus {
		const statuses = [r.blockStatus, r.attestationStatus];
		for (const s of ['missed', 'pending', 'upcoming'] as const) if (statuses.includes(s)) return s;
		return 'done';
	}

	function describe(r: BakerLevelRights) {
		const parts = [];
		if (r.blocks) parts.push(`block: ${r.blockStatus}`);
		if (r.attestations) parts.push(`${r.attestations} attestation slots: ${r.attestationStatus}`);
		return `level ${r.level} — ${parts.join(', ')}`;
	}

	const nowIndex = $derived(rows.findIndex((r) => r.level > level));
</script>

<div class="strip" role="list" aria-label="baking rights around the current level">
	{#each rows as r, i (r.level)}
		{#if i === nowIndex}<span class="now" title="current level {level}"></span>{/if}
		<span role="listitem" class="cell {combined(r)}" class:block={r.blocks > 0} title={describe(r)}></span>
	{/each}
	{#if nowIndex === -1 && rows.length}<span class="now" title="current level {level}"></span>{/if}
	{#if !rows.length}<span class="hint">no rights in the tracked window</span>{/if}
</div>

<style lang="sass">
.strip
	display: flex
	align-items: center
	flex-wrap: wrap
	gap: 4px
	min-height: 16px

.cell
	width: 8px
	height: 8px
	border-radius: 50%
	background-color: var(--success-color)

	&.block
		width: 14px
		height: 14px
		border-radius: 3px

	&.missed
		background-color: var(--error-color)
	&.pending
		background-color: var(--hint-color)
		opacity: 0.6
	&.upcoming
		background-color: transparent
		box-shadow: inset 0 0 0 1.5px var(--accent)

.now
	width: 2px
	height: 18px
	margin: 0 4px
	border-radius: 1px
	background-color: var(--accent)
</style>
