export type TezpayInfo = {
	version: {
		["ami-tezpay"]: string
		tezpay: string
	}
	payout_wallet: string
	configuration: {
		notifications: Array<{ type: string }>
		extensions: Array<{ id: string, kind: string }>
	} & object
}

export const EmptyTezpayInfo: TezpayInfo = {
	version: {
		["ami-tezpay"]: 'unknown',
		tezpay: 'unknown'
	},
	payout_wallet: '',
	configuration: {
		notifications: [],
		extensions: []
	}
}

export type PayoutRecipe = {
	baker: string,
	delegator: string,
	cycle: number,
	recipient: string,
	kind: string,
	fa_token_id: string,
	fa_contract: string,
	delegator_balance: string,
	amount: string,
	fee_rate: number,
	fee: string,
	note: string,
	valid: boolean
}

export type PayoutBlueprintSummary = {
	cycle: number
	delegators : number
	paid_delegators: number
	own_staked_balance: string,
	own_delegated_balance: string,
	external_staked_balance: string,
	external_delegated_balance: string,
	cycle_fees: string,
	cycle_rewards: string,
	distributed_rewards: string,
	bond_income: string,
	fee_income: string,
	total_income: string,
	donated_bonds: string,
	donated_fees: string,
	donated_total: string,
	timestamp: string
}

export type PayoutBlueprint = {
	cycles: number
	payouts: Array<PayoutRecipe>
	summary: PayoutBlueprintSummary
}

export type ExecutionResult = {
	exit_code: number
	success: boolean
}

export type BatchResult = {
	payouts: Array<PayoutRecipe>
	op_hash: string
	is_success: boolean
	err: string | null
}

export type BatchResults = BatchResult[]

export type PayoutReport = {
	id : string
	baker: string
	timestamp: string
	cycle: number
	kind: string
	tx_kind: string 
	contract: string
	token_id: string
	delegator: string
	delegator_balance: string
	recipient: string
	amount: string
	fee_rate: number
	fee: string
	tx_fee: string
	op_hash: string
	success: boolean
	note: string
}

export type CycleReport = {
	name: string
	summary: PayoutBlueprintSummary
	payouts: Array<PayoutReport>
	invalid: Array<PayoutReport>
}