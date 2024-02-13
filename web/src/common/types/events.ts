export type TransactionStage = "building" | "confirming" | "applied" | "failed"

export type TransactionBroadcastedEvent = {
	message: string
	opHash: string
	stage: TransactionStage
}

export type TransactionBuildingEvent = {
	message: string
	stage: TransactionStage
}

export type TransactionErrorEvent = {
	message: string
	error: Error
	stage: TransactionStage
}

export type TransactionEventDispatcher = {
	tx_building: TransactionBuildingEvent
	tx_broadcasted: TransactionBroadcastedEvent
	error: TransactionErrorEvent,
}