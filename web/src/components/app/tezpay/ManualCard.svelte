<script lang="ts">
	import Separator from '../Separator.svelte';
	import Card from '@components/starlight/components/Card.svelte';
	import { createEventDispatcher } from 'svelte';
	import Button from '@components/starlight/components/Button.svelte';
	import Input from '@components/starlight/components/Input.svelte';
	import type { ValidationRules } from '@src/components/starlight/types';
	import CheckBox from '@src/components/starlight/components/CheckBox.svelte';

	const dispatch = createEventDispatcher<{
		pay: { cycle: number | 'latest'; dry: boolean };
	}>();

	let cycle: number | 'latest' = 'latest';
	let dryRun = false;

	const cycleValidationRules: ValidationRules<number | string> = [
		(v) => !!v || "cycle number has to be a number or 'latest'",
		(v) =>
			v === 'latest' ||
			(!isNaN(Number(v)) && Number(v) > 0) ||
			"cycle number has to be a number or 'latest'"
	];

	$: isValidCycle = cycleValidationRules.every((rule) => rule(cycle) === true);
</script>

<div class="manual-wrap">
	<Card class="manual-card">
		<div class="manual">
			<div class="title">
				<h5>Manual Payouts</h5>
			</div>
			<Separator />
			<div class="tools">
				<Input
					label="Cycle"
					hint="Cycle you want to pay out, can be 'latest'"
					rules={cycleValidationRules}
					bind:value={cycle}
				/>
				<div class="dry-run-wrap">
					<div class="dry-run">
						<CheckBox label="Dry Run" bind:checked={dryRun} />
					</div>
				</div>
				<!-- <Button on:click={() => dispatch('preview')}>PREVIEW</Button> -->
			</div>
			<div class="note">NOTE: You will be prompted for a confirmation 😉</div>
			<div class="pay-btn" class:disabled={!isValidCycle}>
				<Button on:click={() => dispatch('pay', { cycle, dry: dryRun })}>Pay</Button>
			</div>
		</div>
	</Card>
</div>

<style lang="sass">
.manual-wrap
	display: grid
	grid-template-rows: 1fr
	width: 100%
	height: 100%
	user-select: none
	// &:hover
	// 	cursor: pointer
	// 	transition: background-color 0.2s
	// 	--card-background-color: #151515

	:global(.manual-card)
		box-sizing: border-box
		height: 100%

.manual
	display: grid
	grid-template-rows: auto auto 1fr
	height: 100%
	gap: var(--spacing)
	.title
		display: flex
		justify-content: center
		h5
			font-size: 1.5rem
			font-weight: 500
			margin: 0

	.tools
		display: flex
		flex-direction: column
		justify-content: center
		align-items: center
		height: 100%
		
		.dry-run-wrap
			display: flex
			align-items: center
			justify-content: end
			width: 100%
			margin-bottom: var(--spacing)
			margin-right: var(--spacing-x2)

	.note
		display: flex
		justify-content: center
		align-items: center
		font-size: 0.9rem

	.pay-btn
		width: 100%
</style>
