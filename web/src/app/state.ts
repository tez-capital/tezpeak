import { derived, writable, type Writable, type Readable } from "svelte/store"

import type { PeakStatus, StatusUpdate } from "@src/common/types"
import { StatusProvider, type StatusProviderStatus } from "./StatusProvider"
import { EMPTY_PEAK_STATUS } from "@src/common/constants"

export const APP_CONNECTION_STATUS = writable("disconnected") as Writable<StatusProviderStatus> 

export const state = writable(EMPTY_PEAK_STATUS) as Writable<PeakStatus>
export const APP_ID = derived(state, $state => {
	const subId = $state?.id
	return subId ? `TEZPEAK - ${subId}` : "TEZPEAK"
})

export const APP_STATUS_LEVEL = derived([state, APP_CONNECTION_STATUS], ([$status, $connectionStatus]) => {
	if ($connectionStatus !== "connected") {
		return "error"
	}

	if ($status?.baker_node.connection_status !== "connected") {
		return "error"
	}

	for (const service of Object.values($status?.services.node_services ?? {})) {
		if (service.status !== "running") {
			return "error"
		}
	}

	for (const service of Object.values($status?.services.signer_services ?? {})) {
		if (service.status !== "running") {
			return "error"
		}
	}

	return "ok"
}) as Readable<"ok" | "error" | "warning">

const provider = new StatusProvider("/sse")
provider.onmessage = (event) => {
	const data = JSON.parse(event.data) as StatusUpdate
	state.set(data.data)
}
provider.onstatuschange = (status) => {
	APP_CONNECTION_STATUS.set(status)
}