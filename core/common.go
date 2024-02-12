package core

import (
	"github.com/tez-capital/tezpeak/core/common"
	"github.com/tez-capital/tezpeak/core/providers"
)

type PeakStatus struct {
	Id        string                          `json:"id,omitempty"` // peak instance id
	BakerNode providers.NodeStatus            `json:"baker_node,omitempty"`
	Nodes     map[string]providers.NodeStatus `json:"nodes,omitempty"`
	Rights    providers.RightsStatus          `json:"rights,omitempty"`
	Services  providers.ServicesStatus        `json:"services,omitempty"`
	Bakers    providers.BakersStatus          `json:"bakers,omitempty"`
	Ledgers   providers.LedgerStatus          `json:"ledgers,omitempty"`
}

type PeakStatusUpdateReportKind string

const (
	FullStatusUpdated    PeakStatusUpdateReportKind = "full"
	PartialStatusUpdated PeakStatusUpdateReportKind = "partial"
)

type PeakStatusUpdatedeRport struct {
	Kind common.StatusUpdateKind `json:"kind,omitempty"`
	Id   string                  `json:"id,omitempty"` // only for partial
	Data PeakStatus              `json:"data,omitempty"`
}
