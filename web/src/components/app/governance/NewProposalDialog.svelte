<script lang="ts">
	import Card from '@src/components/starlight/components/Card.svelte';
	import Input from '@src/components/starlight/components/Input.svelte';
	import Button from '@src/components/starlight/components/Button.svelte';
	import OverlayDialog from '@src/components/starlight/dialogs/Overlay.svelte';
	import { upvote_proposal } from '@src/app/governance';
	import { createEventDispatcher } from 'svelte';
	import type { TransactionEventDispatcher } from '@src/common/types/events';
	import { generateStatementProtocol } from '@src/app/governance/statement';

	export let pkh: string | undefined;
	export let periodIndex: number;
	export let open = false;

	let statement = '';
	let proposal = '';

	$: if (statement) {
		const [ok, proto] = generateStatementProtocol(statement);
		if (ok) {
			proposal = proto;
		}
	}

	$: if (open) {
		statement = '';
		proposal = '';
	}

	$: canInject = proposalValidationRules.every((rule) => rule(proposal) === true);

	const dispatch = createEventDispatcher<TransactionEventDispatcher>();

	async function inject() {
		if (!pkh) return;
		open = false;
		dispatch('tx_building', { message: `building proposal upvote - ${proposal}`, stage: "building" });
		const result = await upvote_proposal(pkh, periodIndex, [proposal]);
		if (typeof result === 'string') {
			dispatch('tx_broadcasted', { message: `proposal upvote - ${proposal}`, opHash: result, stage: "confirming" });
			return;
		}
		dispatch('error', { message: `Failed to upvote. Reason: ${result}`, error: result, stage: "failed"  });
	}

	const proposalValidationRules = [
		(v: string) => !!v || `proposal can not be empty`,
		(v: string) => v.startsWith('P') || `proposal has to start with P`
	];
</script>

<OverlayDialog {open} persistent={true}>
	<div class="center">
		<Card>
			<div class="content">
				<b class="title">New Proposal</b>
				<div />
				<div class="text">
					If you want to make a statement, write it in the <b>statement</b> field; this will generate
					a proposal hash for you.
				</div>
				<div class="text">
					To inject a proper protocol upgrade proposal, write the hash in the <b>new proposal</b>
					field and leave <b>statement</b> blank.
				</div>
				<div class="text">
					*<b>Please refrain from using this feature for spamming.</b>
				</div>
				<div class="text">
					**Remember that you can propose or vote for a proposal only <b>10 times per period</b>.
				</div>

				<div />

				<Input bind:value={statement} label="Statement" />
				<Input
					bind:value={proposal}
					label="New Proposal"
					hint="You can provide proposal hash"
					rules={proposalValidationRules}
				/>

				<div class="controls">
					<Button label="Cancel" on:click={() => (open = false)}></Button>
					<div class:disabled={!canInject}>
						<Button label="Inject" on:click={inject}></Button>
					</div>
				</div>
			</div>
		</Card>
	</div>
</OverlayDialog>

<style lang="sass">
.center
	display: flex
	justify-content: center
	align-items: center
	height: 100%
	width: 100%

.content
	display: grid
	gap: var(--spacing)
	padding: var(--spacing)
	grid-template-columns: 1fr

	.title
		display: flex
		font-size: 1.5rem
		font-weight: 500
		justify-content: center

	.controls
		display: grid
		grid-template-columns: 1fr 1fr
		gap: var(--spacing)

</style>
