import type { PeakStatus } from "./types"

export const TEZPEAK_VERSION = "<VERSION>"
export const TEZPEAK_CODENAME = "<CODENAME>"

export const DEPOSIT_LIMIT_STAKING_CAPACITY_MULTIPLIER = 9n
export const BLOCK_TIME = 15 * 1000 // 15 seconds in milliseconds

export const EMPTY_PEAK_STATUS: PeakStatus = {
	baker_node: {
		connection_status: "disconnected",
	},
	nodes: {

	},
	rights: {
		level: 0,
		rights: [],
	},
	services: {
		timestamp: 0,
		node_services: {},
		signer_services: {},
	},
	bakers: {
		level: 0,
		bakers: {},
	},
	ledgers: {},
}