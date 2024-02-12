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

type Block struct {
	Hash      string    `json:"hash"`
	Timestamp time.Time `json:"timestamp"`
	//	Fitness          string                `json:"fitness"` add if relevant
	Protocol         string                `json:"protocol"`
	LevelInfo        *rpc.LevelInfo        `json:"level_info"`
	VotingPeriodInfo *rpc.VotingPeriodInfo `json:"voting_period_info"`
}

type NodeNetworkInfo struct {
	ConnectionCount int               `json:"connection_count"`
	Stats           *rpc.NetworkStats `json:"stats"`
}

type NodeStatus struct {
	Url              string               `json:"address"`
	ConnectionStatus NodeConnectionStatus `json:"connection_status"`
	Block            *Block               `json:"block"`
	NetworkInfo      *NodeNetworkInfo     `json:"network_info"`
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

type BakerNodeStatusUpdate NodeStatusUpdate

func (s *BakerNodeStatusUpdate) GetId() string {
	return s.Id
}

func (s *BakerNodeStatusUpdate) GetData() interface{} {
	return s.Status
}

func (s *BakerNodeStatusUpdate) GetKind() common.StatusUpdateKind {
	return common.BakerNodeStatusUpdateKind
}

func createMontior(ctx context.Context, client *rpc.Client) (*rpc.BlockHeaderMonitor, error) {
	mon := rpc.NewBlockHeaderMonitor()
	if err := client.MonitorBlockHeader(ctx, mon); err != nil {
		return nil, err
	}
	slog.Debug("created block monitor", "source", client.BaseURL)
	return mon, nil
}

func CollectChainState(ctx context.Context, client *rpc.Client, connectionStateCallback func(state NodeConnectionStatus), blockCallback func(*Block)) {
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

			metadata, err := client.GetBlockMetadata(ctx, h.Hash)
			if err != nil {
				slog.Debug("failed to get block metadata", "source", client.BaseURL.String(), "error", err.Error())
				continue
			}

			blockCallback(&Block{
				Hash:             h.Hash.String(),
				Timestamp:        h.Timestamp,
				Protocol:         metadata.Protocol.String(),
				LevelInfo:        metadata.LevelInfo,
				VotingPeriodInfo: metadata.VotingPeriodInfo,
			})
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
		}, func(h *Block) {
			nodeStatus.Block = h
			statusChannel <- &NodeStatusUpdate{
				Id:     nodeId,
				Status: nodeStatus,
			}
		})
		blockProviderCount.Add(-1)
	}()
}

func updateNetworkInfo(ctx context.Context, client *rpc.Client, nodeStatus *NodeStatus) {
	connections, err := client.GetNetworkConnections(ctx)
	if err == nil {
		nodeStatus.NetworkInfo.ConnectionCount = len(connections)
	} else {
		slog.Debug("failed to get network connections", "source", client.BaseURL.String(), "error", err.Error())
	}

	stats, err := client.GetNetworkStats(ctx)
	if err == nil {
		nodeStatus.NetworkInfo.Stats = stats
	} else {
		slog.Debug("failed to get network stats", "source", client.BaseURL.String(), "error", err.Error())
	}

}

func StartBakerNodeStatusProvider(ctx context.Context, client *rpc.Client, statusChannel chan<- common.ProviderStatusUpdatedReport) {
	go func() {
		nodeStatus := NodeStatus{
			Url:              client.BaseURL.String(),
			ConnectionStatus: Disconnected,
			Block:            nil,
			NetworkInfo: &NodeNetworkInfo{
				ConnectionCount: 0,
				Stats:           nil,
			},
		}

		blockProviderCount.Add(1)
		CollectChainState(ctx, client, func(status NodeConnectionStatus) {
			nodeStatus.ConnectionStatus = status

			updateNetworkInfo(ctx, client, &nodeStatus)

			statusChannel <- &BakerNodeStatusUpdate{
				Status: nodeStatus,
			}
		}, func(h *Block) {
			nodeStatus.Block = h

			updateNetworkInfo(ctx, client, &nodeStatus)

			statusChannel <- &BakerNodeStatusUpdate{
				Status: nodeStatus,
			}
		})
		blockProviderCount.Add(-1)
	}()
}
