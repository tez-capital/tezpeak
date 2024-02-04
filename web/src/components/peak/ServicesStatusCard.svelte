<script lang="ts">
	import Card from '@components/starlight/components/Card.svelte';
	import type { AmiServiceInfo, ServicesStatus } from '@src/common/types';
	import Separator from './Separator.svelte';
	import { formatTimestamp, formatTimestampAgo } from '@src/util/format';
	import { onDestroy } from 'svelte';

	export let title = "Baker's Services";
	export let services: ServicesStatus = {
		timestamp: 0,
		node_services: {},
		signer_services: {}
	};

	type ServiceInfo = AmiServiceInfo & { formattedTimestamp: string };

	$: nodeServices = Object.entries(services.node_services).map((v) => {
		return [v[0], { ...v[1], formattedTimestamp: formatTimestampAgo(v[1].started) }] as [ string, ServiceInfo ];
	});
	$: signerServices = Object.entries(services.signer_services).map((v) => {
		return [v[0], { ...v[1], formattedTimestamp: formatTimestampAgo(v[1].started) } ] as [ string, ServiceInfo ];
	});

	const interval = setInterval(() => {
		for (const [_, v] of nodeServices) {
			v.formattedTimestamp = formatTimestampAgo(v.started);
		}
		for (const [_, v] of signerServices) {
			v.formattedTimestamp = formatTimestampAgo(v.started);
		}
		nodeServices = [...nodeServices]
		signerServices = [...signerServices]
	}, 500);
	onDestroy(() => {
		clearInterval(interval);
	});
</script>

<Card>
	<div class="services">
		<div class="title">
			<h5>{title}</h5>
		</div>

		<div class="services-grid">
			<div />
			<div class="id">Service</div>
			<div class="status">Status</div>
			<div class="active">Active</div>
			<div class="separator">
				<Separator />
			</div>
			{#each nodeServices as [id, serviceInfo]}
				<div class="service-status-color">
					<div
						class="service-status-color-sign"
						class:success={serviceInfo.status === 'running'}
					></div>
				</div>
				<div class="id">{id}</div>
				<div class="status">{serviceInfo.status}</div>
				<div class="active">
					{serviceInfo.status === 'running' ? formatTimestampAgo(serviceInfo.started) : '-'}
				</div>
			{/each}
			<div class="separator">
				<Separator />
			</div>
			{#each signerServices as [id, serviceInfo]}
				<div class="service-status-color">
					<div
						class="service-status-color-sign"
						class:success={serviceInfo.status === 'running'}
					></div>
				</div>
				<div class="id">{id}</div>
				<div class="status">{serviceInfo.status}</div>
				<div class="active">
					{serviceInfo.status === 'running' ? formatTimestampAgo(serviceInfo.started) : '-'}
				</div>
			{/each}
		</div>
		<Separator />
		<div class="timestamp">
			<div class="value">{formatTimestamp(services.timestamp)}</div>
			<!-- <div class="label">updated</div> -->
		</div>
	</div>
</Card>

<style lang="sass">
.services
	display: grid
	gap: var(--spacing)

	.title
		display: flex
		justify-content: center
		h5
			font-size: 1.5rem
			font-weight: 500
			margin: 0
	.services-grid
		display: grid
		gap: var(--spacing)
		grid-template-columns: auto auto auto auto

		.separator
			grid-column: 1/-1

		.id
			text-transform: capitalize
			white-space: nowrap

		.status
			text-transform: uppercase
		
		.active
			display: inline-block
			white-space: nowrap

	.timestamp
		display: grid
		grid-template-columns: auto auto 1fr
		gap: var(--spacing-f2)
		align-items: end

		.value
			grid-column: 1
			font-size: 1.25rem
			font-weight: 500
			display: inline-block
.service-status-color
	.service-status-color-sign
		display: inline-block
		width: 12px
		height: 12px
		border-radius: 50%
		margin-right: var(--spacing-f2)
		background-color: var(--error-color)

		&.success
			background-color: var(--success-color) !important

</style>
