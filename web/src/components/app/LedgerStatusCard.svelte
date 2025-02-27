<script lang="ts">
	import Separator from './Separator.svelte';
	import Card from '@components/starlight/components/Card.svelte';
	import type { LedgerWalletStatus, VotingPeriodInfo } from '@src/common/types/status';

	export let id: string;
	export let info: LedgerWalletStatus;

	$: is_authorized = !!info.authorized;
	$: is_connected = info.ledger_status === 'connected';
	$: app_version = info.app_version ?? "N/A"
</script>

<div class="ledger-wrap">
	<Card class="ledger-card">
		<div class="ledger">
			<div class="title">
				<h5>Ledger Device - {id}</h5>
			</div>
			<Separator />
			<div class="content-grid">
				<div class="sub-title">
					<h6>{info.ledger}</h6>
				</div>
				{#if !is_connected}
					<div class="connection-status error">DISCONNECTED</div>
				{:else if !is_authorized}
					<div class="connection-status error">UNAUTHORIZED</div>
				{:else}
					<div class="connection-status success">READY</div>
				{/if}
				{#if app_version}
					<div class="app-version">
						<h6>App Version: {app_version}</h6>
					</div>
				{/if}
			</div>
		</div>
	</Card>
</div>

<style lang="sass">
.ledger-wrap
	display: grid
	grid-template-rows: 1fr
	width: 100%
	height: 100%
	user-select: none
	// &:hover
	// 	cursor: pointer
	// 	transition: background-color 0.2s
	// 	--card-background-color: #151515

	:global(.ledger-card)
		box-sizing: border-box
		height: 100%

.ledger
	display: grid
	grid-template-rows: auto auto 1fr auto auto
	height: 100%
	gap: var(--spacing)
	.title
		display: flex
		justify-content: center
		h5
			font-size: 1.5rem
			font-weight: 500
			margin: 0
	
	.content-grid
		display: grid
		grid-template-columns: 1fr auto 1fr
		grid-template-rows: auto 3fr auto auto 1fr

		.sub-title
			grid-column: 2
			grid-row: 1
			justify-content: center
			h6
				font-size: 1.1rem
				font-weight: 700
				margin: 0


		.connection-status
			grid-column: 2
			grid-row: 3
			display: flex
			align-items: center
			justify-content: center
			font-size: 1.25rem
			height: 100%

		.app-version
			display: flex
			align-items: center
			justify-content: center
			grid-column: 2
			grid-row: 4
			justify-content: center
			font-size: 1.25rem

.error
	color: var(--error-color)

.success
	color: var(--success-color)
</style>
