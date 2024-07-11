<script lang="ts">
	import { goto } from '$app/navigation';
	import Button from '@components/starlight/components/Button.svelte';
	import AlertDialog from '@components/starlight/dialogs/Alert.svelte';
	import ProgressDialog from '@components/starlight/dialogs/Progress.svelte';

	import ManualCard from '@src/components/app/tezpay/ManualCard.svelte';
	import TestCard from '@src/components/app/tezpay/TestCard.svelte';
	import InfoCard from '@src/components/app/tezpay/InfoCard.svelte';
	import PayoutDialog from '@components/app/tezpay/PayoutDialog.svelte';
	import HomeIcon from '@components/la/icons/home-solid.svelte';
	import ScrollIcon from '@components/la/icons/scroll-solid.svelte';
	import { getTezpayInfo, startContinual, stopContinual } from '@src/app/tezpay/client';

	import { services, wallet } from '@src/app/state/tezpay';
	import { onMount } from 'svelte';
	import { EmptyTezpayInfo, type TezpayInfo } from '@src/common/types/tezpay';

	let info: TezpayInfo = EmptyTezpayInfo;
	let payoutDialog: PayoutDialog;
	let alertDialog: AlertDialog;
	let progressDialog: ProgressDialog;

	onMount(async () => {
		info = await getTezpayInfo();
		payoutDialog.generate('latest');
	});

	async function start() {
		try {
			progressDialog.show({
				title: 'Starting tezpay service...',
				progress: 'indeterminate'
			});
			await startContinual();
		} catch (e: any) {
			await alertDialog.alert({
				title: '⚠️',
				message: e.message
			});
		} finally {
			progressDialog.hide();
		}
	}
	async function stop() {
		try {
			progressDialog.show({
				title: 'Stopping tezpay service...',
				progress: 'indeterminate'
			});
			await stopContinual();
		} catch (e: any) {
			await alertDialog.alert({
				title: '⚠️',
				message: e.message
			});
		} finally {
			progressDialog.hide();
		}
	}
</script>

<div class="dashboard-grid-wrap">
	<div class="navigation-wrap">
		<Button on:click={() => goto('/')}>
			<div class="navigation-btn-content"><HomeIcon /> HOME</div>
		</Button>
		<div />
		<Button on:click={() => goto('/tezpay/reports')}>
			<div class="navigation-btn-content"><ScrollIcon /> REPORTS</div>
		</Button>
	</div>
	<div class="dashboard-grid">
		<div class="info-card">
			<InfoCard
				{info}
				wallet={$wallet}
				services={$services}
				phase="waitinforcycle"
				on:start={start}
				on:stop={stop}
			/>
		</div>
		<!-- <AutomaticCard services={$services}  /> -->
		<ManualCard on:pay={(e) => payoutDialog.generate(e.detail)} />
		<div class="disabled">
			<TestCard />
		</div>
		<!-- <div class="terminal-card">
			<Terminal bind:this={terminal} />
		</div> -->
	</div>
	<PayoutDialog bind:this={payoutDialog} />
	<ProgressDialog bind:this={progressDialog} />
	<AlertDialog bind:this={alertDialog} />
</div>

<style lang="sass">
.dashboard-grid-wrap
	grid-column: 2
	display: grid
	grid-template-columns: 1fr minmax(0px, 1400px) 1fr
	width: calc(100% - var(--spacing) * 2)
	padding: var(--spacing)
	gap: var(--spacing)

	.dashboard-grid
		display: grid
		grid-column: 2
		grid-template-columns: minmax(450px, 1fr) minmax(450px, 1fr) minmax(450px, 1fr)
		gap: var(--spacing)


.navigation-wrap 
	display: grid
	grid-template-columns: auto 1fr auto
	grid-column: 2

	.navigation-btn-content
		display: grid
		grid-template-columns: auto auto
		align-items: center
		gap: var(--spacing)
		:global(svg)
			fill: var(--text-color)
			wdith: 30px
			height: 30px

.info-card
	grid-row: 1/3

.terminal-card
	// width: 200px
	// height: 200px
</style>
