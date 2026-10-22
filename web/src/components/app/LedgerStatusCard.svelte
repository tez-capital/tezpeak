<script lang="ts">
	import { Badge, Card } from '@alis-is/starlight';
	import type { LedgerWalletStatus } from '@src/common/types/status';
	import { formatAddress } from '@src/util/format';

	interface Props {
		id: string;
		info: LedgerWalletStatus;
	}
	let { id, info }: Props = $props();

	const connected = $derived(info.ledger_status === 'connected');
</script>

<Card>
	{#snippet header()}
		<h3 class="sl-card-title">Ledger · {id}</h3>
		{#if !connected}
			<Badge variant="error">disconnected</Badge>
		{:else if !info.authorized}
			<Badge variant="error">unauthorized</Badge>
		{:else}
			<Badge variant="success">ready</Badge>
		{/if}
	{/snippet}

	<dl class="kv">
		<dt>Device</dt>
		<dd>{info.ledger || '-'}</dd>
		<dt>App version</dt>
		<dd>{info.app_version || 'N/A'}</dd>
		{#if info.pkh}
			<dt>Key</dt>
			<dd title={info.pkh}>{formatAddress(info.pkh)}</dd>
		{/if}
	</dl>
</Card>
