export type VotingPeriodInfo = {
	position: number
	remaining: number
	voting_period: {
		index: number
		kind: string
		start_position: number
	}
}

export type NodeStatus = {
	address?: string
	connection_status: "connected" | "disconnected" | "connecting"
	block?: {
		hash: string
		timestamp: string
		//fitness: string
		level_info: {
			level: number
			level_position: number
			cycle: number
			cycle_position: number
		}
		protocol: string
		voting_period_info: VotingPeriodInfo
	}
	network_info?: {
		connection_count: number
		stats?: {
			total_sent: string
			total_recv: string
			current_inflow: number
			current_outflow: number
		} | null
	} | null
}

export type NodesStatus = { [key: string]: NodeStatus }

export type BlockRights = {
	level: number
	rights: { [key: string]: Array<number> }
	realized_checked: boolean
}

export type RightsStatus = {
	level: number
	rights: Array<BlockRights>
	realized_checked?: boolean
}

export type AmiServiceInfo = {
	status: string | "running",
	started: string,
}

export type ServicesStatus = {
	timestamp: number
	node_services: { [key: string]: AmiServiceInfo }
	signer_services: { [key: string]: AmiServiceInfo }
}

export type BakerStatus = {
	deactivated: boolean
	balance: string
	delegated_contracts: Array<string>
	frozen_balance: string
	full_balance: string
	staking_balance: string
	frozen_deposits_limit: string
	delegated_balance: string
	voting_power: string
}

export type BakersStatus = {
	level: number
	bakers: {
		[key: string]: BakerStatus
	}
}

export type LedgerStatus = unknown

export type PeakStatus = {
	id?: string
	baker_node: NodeStatus,
	nodes: NodesStatus
	rights: RightsStatus,
	services: ServicesStatus,
	bakers: BakersStatus,
	ledgers: LedgerStatus,
}

export type StatusUpdate = {
	kind: "full" | "node" | "rights" | "services" | "baker" | "baker-node"
	data: PeakStatus
}

export type NormalizedBlockRights = {
	baker : string,
	blocks: number,
	attestations: number,
	realizedBlocks: number,
	realizedAttestations: number,
}
