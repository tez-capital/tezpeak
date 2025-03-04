import { derived, type Readable } from "svelte/store"
import { state as globalState, nodes } from "."
import type { PeakStatus } from "@src/common/types/status"
import { pickVotingPeriodInfo } from "@src/util/gov"

export const state = derived(globalState, $state => {
	return $state?.modules.tezbake
}) as Readable<PeakStatus["modules"]["tezbake"]>


export const futureBakingRights = derived(state, $tezbakeState => {
	return $tezbakeState?.rights?.rights.filter(right => right.level > $tezbakeState.rights.level).sort((a, b) => a.level - b.level) ?? []
})

export const pastBakingRights = derived(state, $tezbakeState => {
	return $tezbakeState?.rights?.rights.filter(right => right.level <= $tezbakeState.rights.level).sort((a, b) => b.level - a.level) ?? []
})

export const bakers = derived(state, $tezbakeStatus => {
	if ($tezbakeStatus?.bakers === undefined) {
		return []
	}
	return Object.entries($tezbakeStatus.bakers.bakers ?? []).sort(([a], [b]) => a.localeCompare(b))
})

export const wallets = derived(state, $tezbakeStatus => {
	if ($tezbakeStatus?.wallets === undefined) {
		return []
	}
	return Object.entries($tezbakeStatus.wallets ?? []).sort(([a], [b]) => a.localeCompare(b))
})

export const votingPeriodInfo = derived(nodes, $nodes => {
	const nodes = $nodes.map(([, node]) => node)
	return pickVotingPeriodInfo(nodes);
})

export const services = derived(state, $tezbakeStatus => {
	if ($tezbakeStatus === undefined) {
		return { timestamp: 0, applications: {} }
	}

	return $tezbakeStatus.services ?? {}
})

export const status = derived([services, wallets], ([$services, $wallets]) => {
	if ($services === undefined) {
		return "ok"
	}

	for (const [walletId, wallet] of $wallets) {
		if (wallet.ledger_status !== "connected" || !wallet.authorized) {
			return "error"
		}
	}

	// TODO: warnings?
	for (const app of Object.values($services?.applications ?? {})) {
		for (const service of Object.values(app)) {
			if (service.status !== "running") {
				return "error"
			}
		}
	}

	return "ok"
})
