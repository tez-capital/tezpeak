<script lang="ts">
	import { goto } from '$app/navigation';
	import Button from '@components/starlight/components/Button.svelte';
	import AlertDialog from '@components/starlight/dialogs/Alert.svelte';
	import ProgressDialog from '@components/starlight/dialogs/Progress.svelte';
	import SelectDialog from '@components/starlight/dialogs/Select.svelte';
	import TerminalDialog from '@components/starlight/dialogs/Terminal.svelte';
	import ConfirmDialog from '@components/starlight/dialogs/Confirm.svelte';

	import ManualCard from '@src/components/app/tezpay/ManualCard.svelte';
	import TestCard from '@src/components/app/tezpay/TestCard.svelte';
	import InfoCard from '@src/components/app/tezpay/InfoCard.svelte';
	import PayoutDialog from '@components/app/tezpay/PayoutDialog.svelte';
	import HomeIcon from '@components/la/icons/home-solid.svelte';
	import ScrollIcon from '@components/la/icons/scroll-solid.svelte';
	import {
		disableContinual,
		enableContinual,
		getTezpayInfo,
		startContinual,
		stopContinual,
		testExtensions,
		testNotify
	} from '@src/app/tezpay/client';

	import { onMount } from 'svelte';
	import { EmptyTezpayInfo, type TezpayInfo } from '@src/common/types/tezpay';
	import { formatLogMessageForTerminal } from '@src/util/log';
	import { USER_CANCELED } from '@src/components/starlight/src/constants';

	let info: TezpayInfo = EmptyTezpayInfo;
	let payoutDialog: PayoutDialog;
	let alertDialog: AlertDialog;
	let progressDialog: ProgressDialog;
	let selectDialog: SelectDialog;
	let terminalDialog: TerminalDialog;
	let confirmDialog: ConfirmDialog;

	onMount(async () => {
		info = await getTezpayInfo();
		//payoutDialog.generate('latest');
	});

	async function start() {
		try {
			progressDialog.show({
				title: 'Starting continual services...',
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

	async function enable() {
		try {
			progressDialog.show({
				title: 'Enabling continual services...',
				progress: 'indeterminate'
			});
			await enableContinual();
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
				title: 'Stopping continual services...',
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

	async function disable() {
		try {
			await confirmDialog.request_confirmation({
				title: 'Disable Continual Services',
				message: 'Are you sure you want to disable continual payouts?',
				confirmText: 'Yes, Disable',
				cancelText: 'Cancel'
			});

			progressDialog.show({
				title: 'Enabling continual services...',
				progress: 'indeterminate'
			});
			await disableContinual();
		} catch (e: any) {
			if (e === USER_CANCELED) return;

			await alertDialog.alert({
				title: '⚠️',
				message: e.message
			});
		} finally {
			progressDialog.hide();
		}
	}

	async function test_notify() {
		const { configuration } = await getTezpayInfo();
		const notificators = new Set((configuration?.notifications ?? []).map((x) => x.type));

		const result = await selectDialog.prompt({
			title: 'Select Notificator',
			message: 'Select through which notificator you want to send test notifications',
			options: [...notificators, 'all'],
			value: 'all',
			confirmText: 'Send',
			cancelText: 'Cancel'
		});

		terminalDialog.show({
			title: `Test Notifications - ${result}`,
			allowClose: false
		});
		try {
			await testNotify(result, (msg) => {
				terminalDialog.write(formatLogMessageForTerminal(msg));
			});
		} finally {
			terminalDialog.update_state({ allowClose: true });
		}
	}

	async function test_extensions() {
		terminalDialog.show({
			title: `Test Extensions`,
			allowClose: false
		});
		try {
			await testExtensions((msg) => {
				terminalDialog.write(formatLogMessageForTerminal(msg));
			});
		} finally {
			terminalDialog.update_state({ allowClose: true });
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
			<InfoCard phase="" on:start={start} on:stop={stop} on:disable={disable} on:enable={enable} />
		</div>
		<!-- <AutomaticCard services={$services}  /> -->
		<ManualCard on:pay={({ detail: { cycle, dry } }) => payoutDialog.generate(cycle, dry)} />
		<TestCard on:test_extensions={test_extensions} on:test_notifications={test_notify} />

		<!-- <div class="terminal-card">
			<Terminal bind:this={terminal} />
		</div> -->
	</div>
	<PayoutDialog bind:this={payoutDialog} />
	<ProgressDialog bind:this={progressDialog} />
	<AlertDialog bind:this={alertDialog} />
	<SelectDialog bind:this={selectDialog} />
	<TerminalDialog bind:this={terminalDialog} />
	<ConfirmDialog bind:this={confirmDialog} />
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
		grid-template-columns:  repeat(auto-fill, minmax(450px, 1fr))
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
			width: 30px
			height: 30px
</style>
