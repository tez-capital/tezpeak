import { derived, writable, type Writable } from "svelte/store"

import type { PeakStatus } from "@src/common/types/status"
import { StatusProvider, type StatusProviderStatus } from "./provider"
import { EMPTY_PEAK_STATUS } from "@src/common/constants"

export const APP_CONNECTION_STATUS = writable("disconnected") as Writable<StatusProviderStatus>

export const state = writable(EMPTY_PEAK_STATUS) as Writable<PeakStatus>
export const APP_ID = derived(state, $state => {
	const subId = $state?.id
	return subId ? `TEZPEAK - ${subId}` : "TEZPEAK"
})

const provider = new StatusProvider("/api/sse")
provider.onmessage = (event) => {
	const data = JSON.parse(event.data) as PeakStatus
	state.set(data)
}
provider.onstatuschange = (status) => {
	APP_CONNECTION_STATUS.set(status)
}

export const nodes = derived(state, $state => {
	if ($state === undefined) {
		return []
	}
	return Object.entries($state.nodes).sort(([a], [b]) => a.localeCompare(b))
})