package providers

import (
	"context"
	"log/slog"
	"sync/atomic"
	"time"

	"blockwatch.cc/tzgo/rpc"
	"github.com/tez-capital/tezpeak/core/common"
)

type NodeConnectionStatus string

const (
	Connected    NodeConnectionStatus = "connected"
	Connecting   NodeConnectionStatus = "connecting"
	Disconnected NodeConnectionStatus = "disconnected"
)

var (
	blockProviderCount         atomic.Int64
	blockHeaderLogEntryChannel = make(chan *rpc.BlockHeaderLogEntry, 100)
)

type NodeStatus struct {
	Url              string                   `json:"address"`
	ConnectionStatus NodeConnectionStatus     `json:"connection_status"`
	Block            *rpc.BlockHeaderLogEntry `json:"block"`
}

type NodeStatusUpdate struct {
	Id     string
	Status NodeStatus
}

func (s *NodeStatusUpdate) GetId() string {
	return s.Id
}

func (s *NodeStatusUpdate) GetData() interface{} {
	return s.Status
}

func (s *NodeStatusUpdate) GetKind() common.StatusUpdateKind {
	return common.NodeStatusUpdateKind
}

func createMontior(ctx context.Context, client *rpc.Client) (*rpc.BlockHeaderMonitor, error) {
	mon := rpc.NewBlockHeaderMonitor()
	if err := client.MonitorBlockHeader(ctx, mon); err != nil {
		return nil, err
	}
	slog.Debug("created block monitor", "source", client.BaseURL)
	return mon, nil
}

func CollectChainState(ctx context.Context, client *rpc.Client, connectionStateCallback func(state NodeConnectionStatus), blockCallback func(*rpc.BlockHeaderLogEntry)) {
	var err error

	var mon *rpc.BlockHeaderMonitor
	for {
		if mon == nil {
			mon, err = createMontior(ctx, client)
			connectionStateCallback(Connecting)
			if err != nil {
				slog.Debug("failed to create block monitor", "source", client.BaseURL.String(), "error", err.Error())
				connectionStateCallback(Disconnected)
				select {
				case <-ctx.Done():
					return
				case <-time.After(5 * time.Second):
					// Wait 5 seconds before retrying
				}
			} else {
				connectionStateCallback(Connected)
			}
			continue
		}

		select {
		case <-ctx.Done():
			return
		default:
			h, err := mon.Recv(ctx)
			if err == rpc.ErrMonitorClosed {
				slog.Debug("block monitor closed", "source", client.BaseURL.String())
				connectionStateCallback(Disconnected)
				mon = nil
				select {
				case <-ctx.Done():
					return
				case <-time.After(5 * time.Second):
					// Wait 5 seconds before retrying
				}
				continue
			}
			if err != nil && h == nil {
				slog.Debug("failed to receive block", "source", client.BaseURL.String(), "error", err.Error())
				connectionStateCallback(Disconnected)
				mon = nil
				select {
				case <-ctx.Done():
					return
				case <-time.After(5 * time.Second):
					// Wait 5 seconds before retrying
				}
				continue
			}
			go func() {
				blockHeaderLogEntryChannel <- h
			}()
			blockCallback(h)
		}
	}

}

func StartNodeStatusProvider(ctx context.Context, nodeId string, client *rpc.Client, statusChannel chan<- common.ProviderStatusUpdatedReport) {
	go func() {
		nodeStatus := NodeStatus{
			Url:              client.BaseURL.String(),
			ConnectionStatus: Disconnected,
			Block:            nil,
		}
		blockProviderCount.Add(1)
		CollectChainState(ctx, client, func(status NodeConnectionStatus) {
			nodeStatus.ConnectionStatus = status
			statusChannel <- &NodeStatusUpdate{
				Id:     nodeId,
				Status: nodeStatus,
			}
		}, func(h *rpc.BlockHeaderLogEntry) {
			nodeStatus.Block = h
			statusChannel <- &NodeStatusUpdate{
				Id:     nodeId,
				Status: nodeStatus,
			}
		})
		blockProviderCount.Add(-1)
	}()
}
