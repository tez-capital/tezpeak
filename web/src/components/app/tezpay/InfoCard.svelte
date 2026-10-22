<script lang="ts">
	import { Badge, Button, Card } from '@alis-is/starlight';
	import { services, wallet } from '@app/state/tezpay';
	import { formatBalance } from '@src/util/format';
	import { extractContinualServiceInfo } from '@src/util/tezpay';

	interface Props {
		onstart: () => void;
		onstop: () => void;
		onenable: () => void;
		ondisable: () => void;
	}
	let { onstart, onstop, onenable, ondisable }: Props = $props();

	const continual = $derived(extractContinualServiceInfo($services.applications?.tezpay));
	const running = $derived(continual?.status === 'running');
</script>

<Card>
	{#snippet header()}
		<h3 class="sl-card-title">Status</h3>
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
		{#if !continual}
			<Button onclick={onenable}>Enable</Button>
		{:else}
			<div class="danger"><Button onclick={ondisable}>Disable</Button></div>
			{#if running}
				<Button onclick={onstop}>Stop</Button>
			{:else}
				<Button onclick={onstart}>Start</Button>
			{/if}
		{/if}
	{/snippet}
</Card>

<style lang="sass">
.danger
	--button-background-color: transparent
	--button-text-color: var(--error-color)
	--button-hover-background-color: color-mix(in srgb, var(--error-color) 12%, transparent)
	--button-hover-text-color: var(--error-color)
	--button-border: 1px solid color-mix(in srgb, var(--error-color) 45%, transparent)
</style>
