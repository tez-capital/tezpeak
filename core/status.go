package core

import (
	"context"
	"log/slog"

	"blockwatch.cc/tzgo/rpc"
	"github.com/tez-capital/tezpeak/configuration"
	"github.com/tez-capital/tezpeak/constants"
	"github.com/tez-capital/tezpeak/core/common"
	"github.com/tez-capital/tezpeak/core/providers"
)

func Run(ctx context.Context, config *configuration.Runtime) (<-chan PeakStatus, error) {
	status := PeakStatus{
		Nodes:  map[string]providers.NodeStatus{},
		Rights: providers.RightsStatus{},
	}
	statusChannel := make(chan common.ProviderStatusUpdatedReport, 100)

	bakerNodeClient, err := rpc.NewClient(config.Node, nil)
	providers.StartNodeStatusProvider(ctx, constants.BAKER_NODE_ID, bakerNodeClient, statusChannel)
	if err != nil {
		return nil, err
	}

	rightProviderRpcs := make([]*rpc.Client, 0, len(config.ReferenceNodes)+1)
	rightProviderRpcs = append(rightProviderRpcs, bakerNodeClient)

	for id, node := range config.ReferenceNodes {
		if node.Address == "" {
			slog.Warn("no address for node", "id", id)
			continue
		}
		if node.IsRightsProvider {
			if client, err := rpc.NewClient(node.Address, nil); err == nil {
				rightProviderRpcs = append(rightProviderRpcs, client)
			} else {
				slog.Debug("failed to connect to node", "source", node.Address, "error", err.Error())
			}
		}
		if node.IsBlockProvider {
			nodeClient, err := rpc.NewClient(node.Address, nil)
			if err != nil {
				slog.Warn("failed to connect to node", "source", node.Address, "error", err.Error())
				continue
			}
			providers.StartNodeStatusProvider(ctx, id, nodeClient, statusChannel)
		}
	}

	providers.StartRightsStatusProvider(ctx, rightProviderRpcs, config.Bakers, config.BlockWindow, statusChannel)
	// providers.StartServiceStatusProvider(ctx, config.WorkingDirectory, statusChannel)

	resultChannel := make(chan PeakStatus, 100)

	go func() {
		for statusUpdate := range statusChannel {
			switch statusUpdate.GetKind() {
			case common.NodeStatusUpdateKind:
				nodeStatus := statusUpdate.GetData().(providers.NodeStatus)
				status.Nodes[statusUpdate.GetId()] = nodeStatus
			case common.RightsStatusUpdateKind:
				rightsStatus := statusUpdate.GetData().(providers.RightsStatus)
				status.Rights = rightsStatus
			case common.ServicesStatusUpdateKind:
				servicesStatus := statusUpdate.GetData().(providers.ServicesStatus)
				status.Services = servicesStatus
			}
			resultChannel <- status
		}
	}()

	return resultChannel, nil

}
