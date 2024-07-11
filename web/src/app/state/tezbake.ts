import { derived, type Readable } from "svelte/store"
import { state as globalState, nodes } from "."
import type { BlockRights, PeakStatus } from "@src/common/types/status"
import { pickVotingPeriodInfo } from "@src/util/gov"

export const state = derived(globalState, $state => {
	return $state?.modules.tezbake
}) as Readable<PeakStatus["modules"]["tezbake"]>

export const bakingRights = derived(state, $tezbakeState => {
	if ($tezbakeState === undefined) {
		return {
			past: [] as Array<BlockRights>,
			future: [] as Array<BlockRights>,
		}
	}
	return {
		past: $tezbakeState?.rights.rights.filter(right => right.level <= $tezbakeState.rights.level).sort((a, b) => b.level - a.level),
		future: $tezbakeState?.rights.rights.filter(right => right.level > $tezbakeState.rights.level).sort((a, b) => a.level - b.level),
	}
})

export const bakers = derived(state, $tezbakeStatus => {
	if ($tezbakeStatus === undefined) {
		return []
	}
	return Object.entries($tezbakeStatus.bakers.bakers).sort(([a], [b]) => a.localeCompare(b))
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

export const status = derived(services, $services => {
	if ($services === undefined) {
		return "ok"
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
