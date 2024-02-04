import type { PeakStatus } from "./types"

export const PEAK_VERSION = "0.0.1"
export const PEAK_CODENAME = "Vinson"

export const DEPOSIT_LIMIT_STAKING_CAPACITY_MULTIPLIER = 10n
export const BLOCK_TIME = 15 * 1000 // 15 seconds in milliseconds

export const EMPTY_PEAK_STATUS: PeakStatus = {
	// baker_node: {
	// 	connection_status: "disconnected",
	// },
	baker_node: {
		address: "http://localhost:8732",
		connection_status: "connected",
		block: {
			hash: "BMW8LdcE1kGLKpNqsyqGKcsV6Beawiegp9m3HP4bLGFfCPTkRrM",
			timestamp: "2024-02-12T00:40:14Z",
			protocol: "ProxfordYmVfjWnRcgjWH36fW6PArwqykTFzotUxRs6gmTcZDuH",
			level_info: {
				level: 5084150,
				level_position: 5084149,
				cycle: 703,
				cycle_position: 13301,
			},
			voting_period_info: {
				position: 13301,
				remaining: 68618,
				voting_period: {
					index: 116,
					kind: "proposal",
					start_position: 5070848
				}
			}
		}
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