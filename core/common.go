package core

import (
	"encoding/json"
	"log/slog"

	"github.com/tez-capital/tezpeak/core/common"
)

type PeakStatus struct {
	Id      string                       `json:"id,omitempty"` // peak instance id
	Modules map[string]any               `json:"modules,omitempty"`
	Nodes   map[string]common.NodeStatus `json:"nodes,omitempty"`
}

func (s *PeakStatus) MarshalJSON() string {
	resultBytes, err := json.Marshal(s)
	if err != nil {
		slog.Error("failed to marshal peak status", "error", err.Error())
		return "{}"
	}
	return string(resultBytes)
}

type PeakStatusUpdateReportKind string

const (
	FullStatusUpdated    PeakStatusUpdateReportKind = "full"
	PartialStatusUpdated PeakStatusUpdateReportKind = "partial"
)

type PeakStatusUpdatedeRport struct {
	Id   string     `json:"id,omitempty"` // only for partial
	Data PeakStatus `json:"data,omitempty"`
}
