<script lang="ts">
	import { goto } from '$app/navigation';
	import { Button, AlertDialog, ProgressDialog, SelectDialog, ConfirmDialog, USER_CANCELED } from '@alis-is/starlight';
	import TerminalDialog from '@alis-is/starlight/dialogs/Terminal.svelte';

	import ManualCard from '@src/components/app/tezpay/ManualCard.svelte';
	import TestCard from '@src/components/app/tezpay/TestCard.svelte';
	import InfoCard from '@src/components/app/tezpay/InfoCard.svelte';
	import PayoutDialog from '@components/app/tezpay/PayoutDialog.svelte';
	import HomeIcon from '@alis-is/starlight/la/icons/home-solid.svelte';
	import ScrollIcon from '@alis-is/starlight/la/icons/scroll-solid.svelte';
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
				title: 'Disabling continual services...',
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
			allow_close: false
		});
		try {
			await testNotify(result, (msg) => {
				terminalDialog.write(formatLogMessageForTerminal(msg));
			});
		} finally {
			terminalDialog.update_context({ allow_close: true });
		}
	}

	async function test_extensions() {
		terminalDialog.show({
			title: `Test Extensions`,
			allow_close: false
		});
		try {
			await testExtensions((msg) => {
				terminalDialog.write(formatLogMessageForTerminal(msg));
			});
		} finally {
			terminalDialog.update_context({ allow_close: true });
		}
	}
</script>

<div class="tezpay">
	<nav class="page-navigation">
		<Button onclick={() => goto('/')}>
			<div class="navigation-btn-content"><HomeIcon /> HOME</div>
		</Button>
		<div class="spacer"></div>
		<Button onclick={() => goto('/tezpay/reports')}>
			<div class="navigation-btn-content"><ScrollIcon /> REPORTS</div>
		</Button>
	</nav>
	<div class="dashboard-grid">
		<InfoCard onstart={start} onstop={stop} ondisable={disable} onenable={enable} />
		<ManualCard onpay={(cycle, dry) => payoutDialog.generate(cycle, dry)} />
		<TestCard ontestextensions={test_extensions} ontestnotifications={test_notify} />

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
.tezpay
	display: grid
	gap: var(--spacing)

	.dashboard-grid
		display: grid
		grid-template-columns: repeat(auto-fill, minmax(min(100%, 400px), 1fr))
		gap: var(--spacing)
</style>
