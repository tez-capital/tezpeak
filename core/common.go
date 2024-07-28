package core

import (
	"encoding/json"
	"log/slog"
	"sync"

	"github.com/tez-capital/tezpeak/core/common"
)

type peakStatus struct {
	Id        string                     `json:"id,omitempty"` // peak instance id
	Modules   map[string]json.RawMessage `json:"modules,omitempty"`
	Nodes     map[string]json.RawMessage `json:"nodes,omitempty"`
	marshaled []byte                     `json:"-"`

	mtx sync.RWMutex `json:"-"`
}

func newPeakStatus() *peakStatus {
	return &peakStatus{
		Id:      "",
		Modules: make(map[string]json.RawMessage),
		Nodes:   make(map[string]json.RawMessage),
		mtx:     sync.RWMutex{},
	}
}

func (s *peakStatus) SetId(id string) {
	s.Id = id
}

func (s *peakStatus) updateMarshaled() {
	s.marshaled, _ = json.Marshal(s)
}

func (s *peakStatus) UpdateModuleStatus(id string, status any) {
	defer s.updateMarshaled()

	s.mtx.Lock()
	defer s.mtx.Unlock()
	marshaled, err := json.Marshal(status)
	if err != nil {
		slog.Error("failed to marshal module status", "error", err.Error())
		return
	}
	s.Modules[id] = marshaled
}

func (s *peakStatus) UpdateNodeStatus(id string, status common.NodeStatus) {
	defer s.updateMarshaled()

	s.mtx.Lock()
	defer s.mtx.Unlock()
	marshaled, err := json.Marshal(status)
	if err != nil {
		slog.Error("failed to marshal module status", "error", err.Error())
		return
	}
	s.Nodes[id] = marshaled
}

func (s *peakStatus) String() string {
	return string(s.marshaled)
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
