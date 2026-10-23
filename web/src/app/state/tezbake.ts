import { derived, type Readable } from "svelte/store"
import { state as globalState, nodes } from "."
import type { PeakStatus } from "@src/common/types/status"
import { pickVotingPeriodInfo } from "@src/util/gov"

export const state = derived(globalState, $state => {
	return $state?.modules.tezbake
}) as Readable<PeakStatus["modules"]["tezbake"]>


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
