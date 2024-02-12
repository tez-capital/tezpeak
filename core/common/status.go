package common

type StatusUpdateKind string

const (
	FullStatusUpdateKind      StatusUpdateKind = "full"
	NodeStatusUpdateKind      StatusUpdateKind = "node"
	RightsStatusUpdateKind    StatusUpdateKind = "rights"
	ServicesStatusUpdateKind  StatusUpdateKind = "services"
	BakerStatusUpdateKind     StatusUpdateKind = "baker"
	BakerNodeStatusUpdateKind StatusUpdateKind = "baker_node"
	LedgerStatusUpdateKind    StatusUpdateKind = "ledger"
)

type ProviderStatusUpdatedReport interface {
	GetId() string
	GetKind() StatusUpdateKind
	GetData() interface{}
}
