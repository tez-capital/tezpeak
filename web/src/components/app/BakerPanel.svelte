<script lang="ts">
	import { Badge, Card } from '@alis-is/starlight';
	import type { BakerStatus, RightsStatus } from '@src/common/types/status';
	import { bakerRights, getBakerColor, summarizeRights } from '@src/util/baker';
	import { writeToClipboard } from '@src/util/clipboard';
	import { formatBalance, formatBlocksEta, formatTez } from '@src/util/format';
	import { now } from '@app/state/clock';
	import RightsLanes from './RightsLanes.svelte';

	interface Props {
		baker: string;
		status: BakerStatus;
		rights?: RightsStatus;
		showColor?: boolean;
	}
	let { baker, status, rights, showColor = false }: Props = $props();

	// own funds in accent, external in text colour; staked solid, the rest lighter
	const stake = $derived(
		[
			{ label: 'Own staked', mutez: status.staked_balance, color: 'var(--accent)' },
			{
				label: 'Own spendable',
				mutez: status.balance,
				color: 'color-mix(in srgb, var(--accent) 45%, transparent)'
			},
			{
				label: 'External staked',
				mutez: status.external_staked_balance,
				color: 'color-mix(in srgb, var(--text-color) 55%, transparent)'
			},
			{
				label: 'External delegated',
				mutez: status.external_delegated_balance,
				color: 'color-mix(in srgb, var(--text-color) 22%, transparent)'
			}
		].map((s) => ({ ...s, tez: Math.max(0, Number(s.mutez) / 1e6 || 0) }))
	);
	const total = $derived(stake.reduce((sum, s) => sum + s.tez, 0));

	const level = $derived(rights?.level ?? 0);
	const tracked = $derived(rights?.rights ?? []);
	const summary = $derived(summarizeRights(bakerRights(tracked, baker, level)));
	const pastLevels = $derived(tracked.filter((r) => r.level <= level).length);
	const blocks = $derived(summary.blocksDone + summary.blocksMissed);
	const attestations = $derived(summary.attestationsDone + summary.attestationsMissed);
	const eta = $derived(summary.nextBlock && formatBlocksEta(summary.nextBlock.level - level, $now));
</script>

<Card class="baker-panel">
	{#snippet header()}
		<button
			class="unstyle-button baker"
			title="Copy {baker}"
			onclick={() => writeToClipboard(baker)}
		>
			{#if showColor}<span class="swatch" style:background-color={getBakerColor(baker)}></span>{/if}
			<span class="address-text">{baker}</span>
		</button>
		{#if status.deactivated}
			<Badge variant="error">deactivated</Badge>
		{:else}
			<span class="hint nowrap">{status.delegators_count.toLocaleString()} delegators</span>
		{/if}
	{/snippet}

	<!-- the numbers an operator checks first -->
	<div class="figures">
		<div class="figure wide">
			<div class="caption">Staking balance</div>
			<div class="value" title={formatBalance(total * 1e6)}>{formatTez(total * 1e6)}</div>
		</div>
		{#if rights}
			<div class="figure wide">
				<div class="caption">Next block</div>
				{#if summary.nextBlock && eta}
					<div class="value">{eta.in}</div>
					<div class="hint">level {summary.nextBlock.level} · {eta.at}</div>
				{:else}
					<div class="value hint">—</div>
					<div class="hint">none in the next {tracked.length - pastLevels} levels</div>
				{/if}
			</div>
			<div class="figure">
				<div class="caption">Attested</div>
				<div class="value" class:text-error={summary.attestationsMissed > 0}>
					{summary.attestationsDone}/{attestations}
				</div>
				<div class="hint" class:text-error={summary.attestationsMissed > 0}>
					{summary.attestationsMissed
						? `${summary.attestationsMissed} missed`
						: `last ${pastLevels} levels`}
				</div>
			</div>
			<div class="figure">
				<div class="caption">Baked</div>
				<div class="value" class:text-error={summary.blocksMissed > 0}>
					{blocks ? `${summary.blocksDone}/${blocks}` : '—'}
				</div>
				<div class="hint" class:text-error={summary.blocksMissed > 0}>
					{summary.blocksMissed
						? `${summary.blocksMissed} missed`
						: blocks
							? `last ${pastLevels} levels`
							: 'no block rights'}
				</div>
			</div>
		{/if}
	</div>

	{#if rights}
		<section>
			<div class="caption">
				Rights
				<a class="more" href="/rights?baker={baker}">View details →</a>
			</div>
			<RightsLanes rights={tracked} {level} bakers={[baker]} />
		</section>
	{/if}

	<section>
		<div class="caption">Stake composition</div>
		<div class="bar" role="img" aria-label="staking balance composition">
			{#each stake as s (s.label)}
				{#if s.tez > 0}<span
						style:flex={s.tez}
						style:background-color={s.color}
						title="{s.label} {formatTez(s.mutez)}"
					></span>{/if}
			{/each}
		</div>
		<dl class="legend">
			{#each stake as s (s.label)}
				<div>
					<dt><span class="swatch" style:background-color={s.color}></span>{s.label}</dt>
					<dd title={formatBalance(s.mutez)}>{formatTez(s.mutez)}</dd>
				</div>
			{/each}
		</dl>
	</section>
</Card>

<style lang="sass">
.baker
	display: flex
	align-items: center
	gap: var(--spacing-f2)
	min-width: 0
	font-size: 1rem
	letter-spacing: 0.04rem

	.address-text
		overflow: hidden
		text-overflow: ellipsis
		white-space: nowrap

.nowrap
	white-space: nowrap

.swatch
	flex: none
	display: inline-block
	width: 10px
	height: 10px
	border-radius: 3px

:global(.baker-panel .sl-card-body)
	gap: var(--spacing-x2)

section
	display: grid
	gap: var(--spacing-f2)

.caption
	display: flex
	align-items: center
	justify-content: space-between
	font-size: 0.8rem
	letter-spacing: 0.08rem
	text-transform: uppercase
	color: var(--hint-color)

// solid accent pill, like the app's primary buttons, so the way to the rights page is easy to spot
.more
	padding: 6px 16px
	border-radius: 999px
	background-color: var(--accent)
	color: var(--accent-contrast)
	font-size: 0.9rem
	text-transform: none
	letter-spacing: normal
	text-decoration: none
	white-space: nowrap
	transition: background-color 0.15s

	&:hover, &:focus-visible
		background-color: var(--accent-strong)

.figures
	display: grid
	grid-template-columns: repeat(4, minmax(0, 1fr))
	gap: var(--spacing)

.figure
	display: grid
	align-content: start
	gap: 2px
	min-width: 0

	.value
		font-size: 1.6rem
		font-variant-numeric: tabular-nums
		white-space: nowrap
		overflow: hidden
		text-overflow: ellipsis

.bar
	display: flex
	gap: 2px
	height: 10px
	border-radius: 5px
	overflow: hidden
	background-color: var(--progress-bar-background-color)

.legend
	display: grid
	grid-template-columns: repeat(auto-fit, minmax(150px, 1fr))
	gap: var(--spacing-f2) var(--spacing)
	margin: 0

	dt
		display: flex
		align-items: center
		gap: var(--spacing-f2)
		font-size: 0.85rem
		color: var(--hint-color)
		white-space: nowrap

	dd
		margin: 2px 0 0
		padding-left: calc(10px + var(--spacing-f2))
		font-variant-numeric: tabular-nums
		white-space: nowrap

// phones: balance and next block get a full row each, attested/baked share one
@media (max-width: 700px)
	.figures
		grid-template-columns: repeat(2, minmax(0, 1fr))
		row-gap: var(--spacing-x2)

	.figure.wide
		grid-column: 1 / -1

	.legend
		grid-template-columns: repeat(2, minmax(0, 1fr))
</style>
