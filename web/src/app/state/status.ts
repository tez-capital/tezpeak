import { derived, type Readable } from "svelte/store"
import { state, APP_CONNECTION_STATUS } from "."

export type Issue = { level: "error" | "warning", text: string }

// everything that needs the operator's attention; the header status bar and dashboard health line show it
export const issues = derived([state, APP_CONNECTION_STATUS], ([$state, $connectionStatus]) => {
	const result = [] as Array<Issue>
	const error = (text: string) => result.push({ level: "error", text })

	if ($connectionStatus !== "connected") error(`tezpeak ${$connectionStatus}`)

	for (const [id, node] of Object.entries($state.nodes)) {
		if (node.connection_status !== "connected" && node.is_essential) error(`${id} ${node.connection_status}`)
	}

	const { tezbake, tezpay } = $state.modules
	for (const [id, baker] of Object.entries(tezbake?.bakers?.bakers ?? {})) {
		if (baker.deactivated) error(`baker ${id} deactivated`)
	}
	for (const [id, wallet] of Object.entries(tezbake?.wallets ?? {})) {
		if (wallet.ledger_status !== "connected") error(`ledger ${id} ${wallet.ledger_status || "disconnected"}`)
		else if (!wallet.authorized) error(`ledger ${id} not authorized`)
	}
	for (const services of [tezbake?.services, tezpay?.services]) {
		for (const [app, appServices] of Object.entries(services?.applications ?? {})) {
			for (const [id, service] of Object.entries(appServices)) {
				if (service.status !== "running") error(`${app === id ? id : `${app}/${id}`} ${service.status}`)
			}
		}
	}
	if (tezpay?.wallet && tezpay.wallet.level !== "ok") {
		result.push({ level: tezpay.wallet.level === "error" ? "error" : "warning", text: "payout wallet balance low" })
	}
	return result
}) as Readable<Array<Issue>>

export const APP_STATUS_LEVEL = derived(issues, ($issues) =>
	$issues.find((i) => i.level === "error") ? "error" : $issues.length ? "warning" : "ok"
) as Readable<"ok" | "error" | "warning">
