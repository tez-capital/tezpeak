<script lang="ts">
	import { Badge, Card } from '@alis-is/starlight';
	import { services, wallet } from '@app/state/tezpay';
	import { formatBalance, formatTez } from '@src/util/format';
	import { extractContinualServiceInfo } from '@src/util/tezpay';

	const continual = $derived(extractContinualServiceInfo($services.applications?.tezpay));
	const running = $derived(continual?.status === 'running');
</script>

<Card>
	{#snippet header()}
		<a class="title-link" href="/tezpay"
			><h3 class="sl-card-title">Payouts</h3>
			 →</a
		>
		{#if !continual}
			<Badge>manual</Badge>
		{:else if running}
			<Badge variant="success">automatic</Badge>
		{:else}
			<Badge variant="warning">automatic · stopped</Badge>
		{/if}
	{/snippet}

	{#if $wallet !== undefined}
		<div
			class="metric"
			class:text-error={$wallet.level === 'error'}
			class:text-warn={$wallet.level === 'warning'}
			title={formatBalance($wallet.balance)}
		>
			{formatTez($wallet.balance)}
		</div>
		<div class="hint">
			payout wallet{#if $wallet.level !== 'ok'}
				· balance low{/if}
		</div>
	{/if}
</Card>
