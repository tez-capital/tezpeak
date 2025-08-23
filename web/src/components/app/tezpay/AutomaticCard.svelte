<script lang="ts">
	import Card from '@components/starlight/components/Card.svelte';
	import type { AmiServiceInfo, ServicesStatus } from '@src/common/types/status';
	import Separator from '../Separator.svelte';
	import { formatTimestamp, formatTimestampAgo } from '@src/util/format';
	import { createEventDispatcher, onDestroy } from 'svelte';
	import Button from '@components/starlight/components/Button.svelte';
	import { extractContinualServiceInfo } from '@src/util/tezpay';

	export let services: ServicesStatus = {
		timestamp: 0,
		applications: {}
	};

	type ServiceInfo = AmiServiceInfo & { formattedTimestamp: string };

	const interval = setInterval(() => {
		for (const [_, app] of Object.entries(services.applications ?? {})) {
			for (const [_, v] of Object.entries(app)) {
				v.formattedTimestamp = formatTimestampAgo(v.started);
			}
		}
	}, 500);

	$: iterableServices = Object.entries(services.applications ?? {}).map(([id, appServices]) => {
		return [id, Object.entries(appServices)];
	}) as Array<[string, Array<[string, ServiceInfo]>]>;
	onDestroy(() => {
		clearInterval(interval);
	});

	const dispatch = createEventDispatcher<{
		stop: void;
		start: void;
	}>();

	// this is ugly af but it works
	$: hasTezpayStatus = !!extractContinualServiceInfo(services.applications?.tezpay);
	$: isTezpayRunning = extractContinualServiceInfo(services.applications?.tezpay)?.status === 'running';
</script>

<Card>
	<div class="services">
		<div class="title">
			<h5>Automatic Payouts</h5>
		</div>
		<div class="separator">
			<Separator />
		</div>

		<div class="services-grid">
			<div />
			<div class="id">Service</div>
			<div class="status">Status</div>
			<div class="active">Active</div>

			<div class="separator-inner">
				<Separator />
			</div>

			{#each iterableServices as [id, appServices]}
				<div class="application-title">{id}</div>

				{#if appServices.length === 0}
					<div class="no-services">N/A</div>
				{:else}
					{#each appServices as [id, serviceInfo]}
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
				{/if}
			{/each}
		</div>
		<div class="timestamp">
			<div class="value">{formatTimestamp(services.timestamp)}</div>
			<div class="separator">
				<Separator />
			</div>
			<!-- <div class="label">updated</div> -->
		</div>

		{#if hasTezpayStatus}
			<div class="tools">
				{#if isTezpayRunning}
					<Button on:click={() => dispatch('stop')}>STOP</Button>
				{:else}
					<Button on:click={() => dispatch('start')}>START</Button>
				{/if}
			</div>
		{/if}
	</div>
</Card>

<style lang="sass">
.services
	display: grid
	gap: var(--spacing)
	grid-template-rows: auto auto 1fr auto
	height: 100%

	.title
		display: flex
		justify-content: center
		margin-bottom: var(--spacing)
		h5
			font-size: 1.5rem
			font-weight: 500
			margin: 0
	.services-grid
		display: grid
		gap: var(--spacing)
		grid-template-columns: auto auto auto auto
		margin: var(--spacing)

		.id
			text-transform: capitalize
			white-space: nowrap

		.status
			text-transform: uppercase
		
		.active
			display: inline-block
			white-space: nowrap

		.separator-inner
			grid-column: 1/-1
			justify-self: center
			width: 90%
			opacity: 0.5

	.timestamp
		grid-column: 1/-1
		width: 100%
		display: grid
		grid-template-columns: 1fr auto  
		gap: var(--spacing-f2)
		align-items: end

		.value
			grid-column: 2
			font-size: .9rem
			font-weight: 500
			display: flex

			
.separator
	grid-column: 1/-1

.application-title
	grid-column: 1/-1
	font-weight: bold
	text-transform: capitalize
	text-align: center

.no-services
	grid-column: 1/-1
	text-align: center

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
