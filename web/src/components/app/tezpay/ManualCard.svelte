<script lang="ts">
	import { Button, Card, CheckBox, Input, type ValidationRules } from '@alis-is/starlight';

	interface Props {
		onpay: (cycle: number | 'latest', dry: boolean) => void;
	}
	let { onpay }: Props = $props();

	let cycle = $state<number | 'latest'>('latest');
	let dryRun = $state(false);

	const cycleValidationRules: ValidationRules<number | string> = [
		(v) => !!v || "cycle number has to be a number or 'latest'",
		(v) =>
			v === 'latest' ||
			(!isNaN(Number(v)) && Number(v) > 0) ||
			"cycle number has to be a number or 'latest'"
	];

	const isValidCycle = $derived(cycleValidationRules.every((rule) => rule(cycle) === true));
</script>

<Card>
	{#snippet header()}
		<h3 class="sl-card-title">Manual payouts</h3>
	{/snippet}

	<Input
		label="Cycle"
		hint="Cycle you want to pay out, can be 'latest'"
		rules={cycleValidationRules}
		bind:value={cycle}
	/>
	<CheckBox label="Dry run" bind:checked={dryRun} />

	{#snippet footer()}
		<span class="hint">You review the payouts before anything is sent.</span>
		<div class:disabled={!isValidCycle}>
			<Button onclick={() => onpay(cycle, dryRun)}>Pay</Button>
		</div>
	{/snippet}
</Card>
