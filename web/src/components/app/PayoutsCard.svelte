<script lang="ts">
	import { goto } from '$app/navigation';
	import { Badge, Button, Card } from '@alis-is/starlight';
	import { services, wallet } from '@app/state/tezpay';
	import { formatBalance } from '@src/util/format';
	import { extractContinualServiceInfo } from '@src/util/tezpay';

	const continual = $derived(extractContinualServiceInfo($services.applications?.tezpay));
	const running = $derived(continual?.status === 'running');
</script>

<Card>
	{#snippet header()}
		<h3 class="sl-card-title">Payouts</h3>
		{#if !continual}
			<Badge>disabled</Badge>
		{:else if running}
			<Badge variant="success">active</Badge>
		{:else}
			<Badge variant="warning">inactive</Badge>
		{/if}
	{/snippet}

	<dl class="kv">
		<dt>Automatic payouts</dt>
		<dd class:text-ok={running} class:text-warn={continual && !running}>
			{!continual ? 'disabled' : running ? 'active' : 'inactive'}
		</dd>
		{#if $wallet !== undefined}
			<dt>Wallet balance</dt>
			<dd class:text-error={$wallet.level === 'error'} class:text-warn={$wallet.level === 'warning'}>
				{formatBalance($wallet.balance)}
			</dd>
		{/if}
	</dl>

	{#snippet footer()}
		<Button onclick={() => goto('/tezpay')}>Open</Button>
	{/snippet}
</Card>
