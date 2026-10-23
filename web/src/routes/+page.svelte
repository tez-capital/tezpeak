<script lang="ts">
	import { nodes } from '@app/state/index';
	import { issues } from '@app/state/status';
	import {
		state as tezbakeStatus,
		bakers as tezbakeBakers,
		services as tezbakeServices,
		wallets as tezbakeWallets,
		votingPeriodInfo
	} from '@app/state/tezbake';
	import { state as tezpayStatus } from '@app/state/tezpay';
	import { now } from '@app/state/clock';
	import BakerPanel from '@components/app/BakerPanel.svelte';
	import InfrastructureCard from '@components/app/InfrastructureCard.svelte';
	import GovernancePeriodCard from '@components/app/GovernancePeriodCard.svelte';
	import PayoutsCard from '@components/app/PayoutsCard.svelte';
	import { formatTimestampAgoStrict } from '@src/util/format';

	// chain head from the most advanced connected node
	const head = $derived(
		$nodes
			.map(([, n]) => n.block)
			.filter((b) => b !== undefined)
			.sort((a, b) => b.level_info.level - a.level_info.level)[0]
	);
	const errors = $derived($issues.filter((i) => i.level === 'error').length);
</script>

<div class="dashboard">
	<section class="overview" class:error={errors > 0} class:warning={!errors && $issues.length > 0}>
		<div class="health">
			<span
				class="status-dot"
				class:ok={!$issues.length}
				class:error={errors > 0}
				class:warn={!errors && $issues.length > 0}
			></span>
			{#if !$issues.length}
				All systems operational
			{:else}
				<div>
					{$issues.length}
					{$issues.length === 1 ? 'issue needs' : 'issues need'} attention
					<ul>
						{#each $issues as issue, i (i)}
							<li
								class:text-error={issue.level === 'error'}
								class:text-warn={issue.level === 'warning'}
							>
								{issue.text}
							</li>
						{/each}
					</ul>
				</div>
			{/if}
		</div>
		{#if head}
			<dl class="chain">
				<div>
					<dt>Level</dt>
					<dd>{head.level_info.level}</dd>
				</div>
				<div>
					<dt>Cycle</dt>
					<dd>{head.level_info.cycle}</dd>
				</div>
				<div>
					<dt>Last block</dt>
					<dd>
						{#key $now}{formatTimestampAgoStrict(head.timestamp)}{/key}
					</dd>
				</div>
			</dl>
		{/if}
	</section>

	{#if $tezbakeStatus}
		{#each $tezbakeBakers as [baker, info] (baker)}
			<BakerPanel
				{baker}
				status={info}
				rights={$tezbakeStatus.rights}
				showColor={$tezbakeBakers.length > 1}
			/>
		{:else}
			<p class="hint">No bakers are tracked.</p>
		{/each}
	{/if}

	<div class="cards">
		{#if $tezpayStatus}<PayoutsCard />{/if}
		{#if $tezbakeStatus}<GovernancePeriodCard votingPeriodInfo={$votingPeriodInfo} />{/if}
		<InfrastructureCard nodes={$nodes} services={$tezbakeServices} wallets={$tezbakeWallets} />
	</div>
</div>

<style lang="sass">
.dashboard
	display: grid
	grid-template-columns: minmax(0, 1fr)
	gap: var(--spacing)

.overview
	display: flex
	flex-wrap: wrap
	align-items: center
	justify-content: space-between
	gap: var(--spacing)
	padding: var(--spacing) var(--spacing-x2)
	border-radius: var(--border-radius)
	border: 1px solid var(--frost-border-color)
	border-left: 3px solid var(--success-color)
	background-color: var(--card-background-color)

	&.warning
		border-left-color: var(--warning-color)
	&.error
		border-left-color: var(--error-color)

.health
	display: flex
	align-items: baseline
	gap: var(--spacing)
	font-size: 1.1rem

	ul
		margin: 4px 0 0
		padding: 0
		list-style: none
		display: flex
		flex-wrap: wrap
		gap: 2px var(--spacing)
		font-size: 0.9rem

.chain
	display: flex
	gap: var(--spacing-x2)
	margin: 0

	dt
		font-size: 0.75rem
		letter-spacing: 0.08rem
		text-transform: uppercase
		color: var(--hint-color)

	dd
		margin: 0
		font-size: 1.1rem
		font-variant-numeric: tabular-nums

// payouts, governance, infrastructure side by side under the bakers
.cards
	display: grid
	grid-template-columns: repeat(auto-fit, minmax(min(100%, 320px), 1fr))
	align-items: start
	gap: var(--spacing)
</style>
