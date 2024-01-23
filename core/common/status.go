package common

type StatusUpdateKind string

const (
	FullStatusUpdateKind     StatusUpdateKind = "full"
	NodeStatusUpdateKind     StatusUpdateKind = "node"
	RightsStatusUpdateKind   StatusUpdateKind = "rights"
	ServicesStatusUpdateKind StatusUpdateKind = "services"
)

type ProviderStatusUpdatedReport interface {
	GetId() string
	GetKind() StatusUpdateKind
	GetData() interface{}
}
