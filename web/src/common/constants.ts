import type { PeakStatus } from "./types/status"

export const TEZPEAK_VERSION = "<VERSION>"
export const TEZPEAK_CODENAME = "<CODENAME>"

export const DEPOSIT_LIMIT_DELEGATION_CAPACITY_MULTIPLIER = 9n
export const BLOCK_TIME = 6 * 1000 // 6 seconds in milliseconds

export const EXPLORER_URL = "https://tzkt.io"

// governance
export const PROPOSAL_PERIOD_QUORUM_PERCENTAGE = 5n
export const PROPOSAL_REQUIRED_MAJORITY_PERCENTAGE = 80n
export const AGORA_PROPOSAL_URL = "https://www.tezosagora.org/proposal"

export const EMPTY_PEAK_STATUS: PeakStatus = {
	modules: {
		tezbake: undefined,
		tezpay: undefined
	},
	nodes: {}
}