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

const provider = new StatusProvider("/sse")
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
/*

	logger.Info("dry running batch", "id", batchId, "tx_count", len(batch))
	if state.Global.GetWantsOutputJson() {
		logger.Info("creating batch", "recipes", batch, "phase", "executing_batch")
	} else {
		logger.Info("creating batch", "tx_count", len(batch), "phase", "executing_batch")
	}
	opExecCtx, err := batch.ToOpExecutionContext(ctx.GetSigner(), ctx.GetTransactor())
	if err != nil {
		logger.Warn("failed to create operation execution context", "id", batchId, "error", err.Error(), "phase", "batch_execution_finished")
		return common.NewFailedBatchResultWithOpHash(batch, opExecCtx.GetOpHash(), errors.Join(constants.ErrOperationContextCreationFailed, err))
	}
	logger.Info("broadcasting batch")
	time.Sleep(2 * time.Second)
	logger.Info("waiting for confirmation", "op_reference", utils.GetOpReference(opExecCtx.GetOpHash(), ctx.GetConfiguration().Network.Explorer), "op_hash", opExecCtx.GetOpHash(), "phase", "batch_waiting_for_confirmation")
	time.Sleep(4 * time.Second)
	logger.Info("batch successful", "phase", "batch_execution_finished")
	return common.NewSuccessBatchResult(batch, tezos.ZeroOpHash)
*/