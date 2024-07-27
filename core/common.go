package core

import (
	"encoding/json"
	"log/slog"
	"sync"

	"github.com/tez-capital/tezpeak/core/common"
)

type peakStatus struct {
	Id      string                       `json:"id,omitempty"` // peak instance id
	Modules map[string]any               `json:"modules,omitempty"`
	Nodes   map[string]common.NodeStatus `json:"nodes,omitempty"`

	mtx sync.RWMutex
}

func newPeakStatus() *peakStatus {
	return &peakStatus{
		Id:      "",
		Modules: make(map[string]any),
		Nodes:   make(map[string]common.NodeStatus),
		mtx:     sync.RWMutex{},
	}
}

func (s *peakStatus) SetId(id string) {
	s.Id = id
}

func (s *peakStatus) UpdateModuleStatus(id string, status any) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.Modules == nil {
		s.Modules = make(map[string]any)
	}
	s.Modules[id] = status
}

func (s *peakStatus) UpdateNodeStatus(id string, status common.NodeStatus) {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	if s.Nodes == nil {
		s.Nodes = make(map[string]common.NodeStatus)
	}
	s.Nodes[id] = status
}

func (s *peakStatus) ToJSONString() string {
	s.mtx.RLock()
	defer s.mtx.RUnlock()

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
	Data peakStatus `json:"data,omitempty"`
}
