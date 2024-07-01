package core

import (
	"context"
	"log/slog"

	"github.com/tez-capital/tezbake/apps/base"
	"github.com/tez-capital/tezpeak/configuration"
	"github.com/tez-capital/tezpeak/constants/enums"
	"github.com/tez-capital/tezpeak/core/common"
	"github.com/tez-capital/tezpeak/core/providers"
	"github.com/trilitech/tzgo/rpc"
)

func Run(ctx context.Context, config *configuration.Runtime) (<-chan PeakStatus, error) {
	status := PeakStatus{
		Id:    config.Id,
		Nodes: map[string]providers.NodeStatus{},
		Rights: providers.RightsStatus{
			Level:  0,
			Rights: []*providers.BlockRights{},
		},
		Services: providers.ServicesStatus{
			NodeServices:   map[string]base.AmiServiceInfo{},
			SignerServices: map[string]base.AmiServiceInfo{},
		},
		Bakers: providers.BakersStatus{
			Level:  0,
			Bakers: map[string]*providers.BakerStakingStatus{},
		},
		Ledgers: providers.LedgerStatus{
			Level: 0,
		},
	}
	statusChannel := make(chan common.ProviderStatusUpdatedReport, 100)

	bakerNodeClient, err := rpc.NewClient(config.Node, nil)
	providers.StartBakerNodeStatusProvider(ctx, bakerNodeClient, statusChannel)
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
	providers.StartBakersStatusProvider(ctx, rightProviderRpcs, config.Bakers, statusChannel)
	switch config.Providers.Services {
	case enums.TezbakeServiceStatusProvider:
		providers.StartServiceStatusProvider(ctx, config.TezbakeHome, statusChannel)
	case enums.NoneServiceStatusProvider:
	default:
		slog.Warn("unknown service status provider", "provider", config.Providers.Services)
	}

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
			case common.BakerStatusUpdateKind:
				bakersStatus := statusUpdate.GetData().(providers.BakersStatus)
				status.Bakers = bakersStatus
			case common.LedgerStatusUpdateKind:

			case common.BakerNodeStatusUpdateKind:
				nodeStatus := statusUpdate.GetData().(providers.NodeStatus)
				status.BakerNode = nodeStatus
			}
			resultChannel <- status
		}
	}()

	return resultChannel, nil

}
