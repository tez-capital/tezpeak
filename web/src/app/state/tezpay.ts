import { derived, type Readable } from "svelte/store"
import { state as globalState } from "."
import type { PeakStatus } from "@src/common/types/status"

export const state = derived(globalState, $state => {
	return $state?.modules.tezpay
}) as Readable<PeakStatus["modules"]["tezpay"]>

export const services = derived(state, $state => {
	if ($state === undefined) {
		return { timestamp: 0, applications: {} }
	}

	return $state.services ?? {}
})

export const wallet = derived(state, $state => {
	if ($state === undefined) {
		return undefined
	}

	return $state.wallet ?? { address: "", balance: 0, level: "ok" }
})

export const status = derived([services, wallet], ([$services, $wallet]) => {
	// TODO: warnings?
	if ($services !== undefined) {
		for (const app of Object.values($services.applications ?? {})) {
			for (const service of Object.values(app)) {
				if (service.status !== "running") {
					return "error"
				}
			}
		}
	}

	if ($wallet !== undefined && $wallet.level !== "ok") {
		return $wallet.level
	}


	return "ok"
})

