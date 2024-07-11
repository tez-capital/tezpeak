import { derived, type Readable } from "svelte/store"
import { state, APP_CONNECTION_STATUS } from "."
import { status as tezbakeStatus } from "./tezbake"
import { status as tezpayStatus } from "./tezpay"

export const APP_STATUS_LEVEL = derived([state, APP_CONNECTION_STATUS, tezbakeStatus, tezpayStatus], ([$state, $connectionStatus, $tezbakeStatus, $tezpayStatus]) => {
	if ($connectionStatus !== "connected") {
		return "error"
	}

	for (const node of Object.values($state.nodes)) {
		if (node.connection_status !== "connected" && node.is_essential) {
			return "error"
		}
	}
	if ($tezbakeStatus !== "ok") {
		return $tezbakeStatus
	}

	if ($tezpayStatus !== "ok") {
		return $tezpayStatus
	}

	return "ok"
}) as Readable<"ok" | "error" | "warning">