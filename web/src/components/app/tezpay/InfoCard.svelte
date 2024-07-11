<script lang="ts">
	import Separator from '../Separator.svelte';
	import Card from '@components/starlight/components/Card.svelte';
	import { createEventDispatcher, onDestroy } from 'svelte';
	import Button from '../../starlight/components/Button.svelte';
	import type { TezpayInfo } from '@src/common/types/tezpay';
	import type { ServicesStatus, WalletStatus } from '@src/common/types/status';
	import { formatAddress, formatBalance, formatTimestampAgo } from '@src/util/format';
	import { writeToClipboard } from '@src/util/clipboard';

	export let info: TezpayInfo;
	export let wallet: WalletStatus | undefined;
	export let services: ServicesStatus = {
		timestamp: 0,
		applications: {}
	};
	export let phase: string;

	// TODO: new version warning

	const dispatch = createEventDispatcher<{
		view_configuration: void;
		stop: void;
		start: void;
	}>();

	$: appServices = Object.entries(services.applications?.tezpay ?? {});
	$: isTezpayRunning = services.applications?.tezpay?.['tezpay']?.status === 'running';
	$: hasTezpayStatus = !!services.applications?.tezpay?.['tezpay'];
</script>

<div class="governance-wrap">
	<Card class="governance-card">
		<div class="governance">
			<div class="title">
				<h5>TEZPAY</h5>
			</div>
			<Separator />
			<div class="info-data-grid">
				<div class="subtitle">Versions</div>
				<div class="package">Package</div>
				<div class="version">Version</div>
				<div class="package">ami-tezpay</div>
				<div class="version">{info.version['ami-tezpay']}</div>

				<div class="package">tezpay</div>
				<div class="version">{info.version.tezpay}</div>
				<div class="separator">
					<Separator />
				</div>
				{#if wallet}
					<div class="subtitle">Balance</div>
					<button class="unstyle-button wallet" on:click={() => writeToClipboard(wallet.address)}
						>{formatAddress(wallet.address)}</button
					>
					<div
						class="balance"
						class:error={wallet.level === 'error'}
						class:warn={wallet.level === 'warning'}
					>
						{formatBalance(wallet.balance)}
					</div>
				{/if}
				{#if appServices.length > 0}
					<div class="separator">
						<Separator />
					</div>
					<div class="subtitle">Services</div>
					{#each appServices as [id, serviceInfo]}
						<div class="services-grid">
							<div />
							<div class="id">Service</div>
							<div class="status">Status</div>
							<div class="active">Active</div>
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
						</div>
					{/each}
					<div class="separator">
						<Separator />
					</div>
					<div class="subtitle">Automatic Payouts</div>
					{#if isTezpayRunning}
						<div class="automatic-payouts-status ok">ACTIVE</div>
						<div class="automatic-payouts-phase">{phase}</div>
					{:else}
						<div class="automatic-payouts-status warn">INACTIVE</div>
						<div class="automatic-payouts-phase"></div>
					{/if}
					{#if hasTezpayStatus}
						<div class="tools">
							{#if isTezpayRunning}
								<Button on:click={() => dispatch('stop')}>STOP</Button>
							{:else}
								<Button on:click={() => dispatch('start')}>START</Button>
							{/if}
						</div>
						<div class="separator">
							<Separator />
						</div>
					{/if}
				{/if}
			</div>
			<!-- <Separator />
			<Button on:click={() => dispatch('view_configuration')}>View Configuration</Button> -->
		</div>
	</Card>
</div>

<style lang="sass">
.governance-wrap
	display: grid
	grid-template-rows: 1fr
	width: 100%
	height: 100%
	user-select: none
	// &:hover
	// 	cursor: pointer
	// 	transition: background-color 0.2s
	// 	--card-background-color: #151515

	:global(.governance-card)
		box-sizing: border-box
		height: 100%

.governance
	display: grid
	grid-template-rows: auto auto auto auto 1fr
	height: 100%
	gap: var(--spacing)
	.title
		display: flex
		justify-content: center
		h5
			font-size: 1.5rem
			font-weight: 500
			margin: 0

.info-data-grid
	display: grid
	gap: var(--spacing)
	grid-template-columns: auto auto 
	margin: var(--spacing)

	.subtitle
		font-weight: 500
		justify-self: center
		text-transform: uppercase
		font-size: 0.9rem
		grid-column: 1/-1

	.package
		white-space: nowrap

	.version
		font-size: 0.9rem
		white-space: nowrap

.separator
	grid-column: 1/-1
	margin: var(--spacing) 0

.error 
	color: var(--error-color)

.warn
	color: var(--warning-color)

.ok
	color: var(--success-color)

.services-grid
	grid-column: 1/-1
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

.automatic-payouts-status
	grid-column: 1/-1
	height: 100px
	display: flex
	justify-content: center
	align-items: center

.tools
	grid-column: 1/-1

.automatic-payouts-phase
	grid-column: 1/-1
	font-size: 0.9rem
	display: flex
	justify-content: center
	align-items: center
</style>
