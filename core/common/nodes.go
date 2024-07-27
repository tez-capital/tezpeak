package common

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/tez-capital/tezpeak/configuration"
	"github.com/tez-capital/tezpeak/constants"
	"github.com/trilitech/tzgo/rpc"
)

type NodeNetworkInfo struct {
	ConnectionCount int               `json:"connection_count"`
	Stats           *rpc.NetworkStats `json:"stats"`
}

type NodeStatus struct {
	Url              string           `json:"address"`
	ConnectionStatus ConnectionStatus `json:"connection_status"`
	Block            *Block           `json:"block"`
	NetworkInfo      *NodeNetworkInfo `json:"network_info"`
	IsEssential      bool             `json:"is_essential"`
}

type NodeStatusUpdate struct {
	Id     string
	Status NodeStatus
}

func (s *NodeStatusUpdate) GetId() string {
	return s.Id
}

func (s *NodeStatusUpdate) GetData() any {
	return s.Status
}

type ActiveRpcNode struct {
	configuration.TezosNode
	*rpc.Client
}

var (
	activeRpcNodes    = make(map[string]*ActiveRpcNode)
	defaultHttpClient = &http.Client{
		Timeout: constants.DEFAULT_HTTP_TIMEOUT_SECONDS * time.Second,
	}
)

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

func StartNodeStatusProviders(ctx context.Context, nodes map[string]configuration.TezosNode, statusChannel chan<- StatusUpdate) {
	for nodeId, node := range nodes {

		if _, ok := activeRpcNodes[nodeId]; ok {
			slog.Warn("node already active", "source", node.Address, "id", nodeId)
			continue
		}
		client, err := rpc.NewClient(node.Address, defaultHttpClient)
		if err != nil {
			slog.Warn("failed to connect to node", "source", node.Address, "error", err.Error())
			continue
		}

		activeRpcNodes[nodeId] = &ActiveRpcNode{
			TezosNode: node,
			Client:    client,
		}

		if !node.IsBlockProvider {
			continue
		}
		// default blockMonitorClient uses timeout, but for streaming we do not want one so use nil instead
		blockMonitorClient, err := rpc.NewClient(node.Address, nil)
		if err != nil {
			slog.Warn("failed to connect to node", "source", node.Address, "error", err.Error())
			continue
		}

		nodeStatus := NodeStatus{
			Url:              node.Address,
			ConnectionStatus: Disconnected,
			Block:            nil,
			NetworkInfo:      nil,
			IsEssential:      node.IsEssential,
		}

		if node.IsNetworkInfoProvider {
			nodeStatus.NetworkInfo = &NodeNetworkInfo{}
		}

		monitorId, err := AddBlockMonitor(ctx, blockMonitorClient, func(status ConnectionStatus) {
			nodeStatus.ConnectionStatus = status

			if node.IsNetworkInfoProvider {
				updateNetworkInfo(ctx, blockMonitorClient, &nodeStatus)
			}

			statusChannel <- &NodeStatusUpdate{
				Id:     nodeId,
				Status: nodeStatus,
			}
		}, func(h *Block) {
			nodeStatus.Block = h

			if node.IsNetworkInfoProvider {
				updateNetworkInfo(ctx, blockMonitorClient, &nodeStatus)
			}

			statusChannel <- &NodeStatusUpdate{
				Id:     nodeId,
				Status: nodeStatus,
			}
		})
		if err != nil {
			slog.Warn("failed to add block monitor", "source", blockMonitorClient.BaseURL.String(), "error", err.Error())
			continue
		}

		go func() {
			<-ctx.Done()
			RemoveBlockMonitor(monitorId)
		}()
	}
}

func isClientSynced(ctx context.Context, client *rpc.Client) bool {
	status, err := client.GetStatus(ctx)
	return status.SyncState == "synced" || (err != nil && strings.Contains(err.Error(), "status 403"))
}

func AttemptWithRpcClients[T any](ctx context.Context, f func(client *ActiveRpcNode) (T, error)) (T, error) {
	var err error
	var result T

	for _, node := range activeRpcNodes {
		if !isClientSynced(ctx, node.Client) {
			continue
		}
		slog.Debug("attempting with client", "client", node.Client.BaseURL.Host)

		result, err = f(node)
		if err != nil {
			continue
		}
		return result, nil
	}
	return result, err
}
