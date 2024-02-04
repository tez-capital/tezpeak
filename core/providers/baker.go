package providers

import (
	"context"
	"log/slog"
	"strings"

	"blockwatch.cc/tzgo/rpc"
	"blockwatch.cc/tzgo/tezos"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/tez-capital/tezpeak/core/common"
)

type BakersStatus struct {
	Level  int64                    `json:"level"`
	Bakers map[string]*rpc.Delegate `json:"bakers"`
}

type BakersStatusUpdate struct {
	Status BakersStatus
}

func (s *BakersStatusUpdate) GetId() string {
	return "bakers"
}

func (s *BakersStatusUpdate) GetData() interface{} {
	return s.Status
}

func (s *BakersStatusUpdate) GetKind() common.StatusUpdateKind {
	return common.BakerStatusUpdateKind
}

func getBakerStatusFor(ctx context.Context, clients []*rpc.Client, baker string) (*rpc.Delegate, error) {
	addr, err := tezos.ParseAddress(baker)
	if err != nil {
		return nil, err
	}
	status, err := attemptWithClients(clients, func(client *rpc.Client) (*rpc.Delegate, error) {
		acc, err := client.GetDelegate(ctx, addr, rpc.Head)
		if err != nil {
			return nil, err
		}
		return acc, nil
	})
	return status, err
}

func StartBakersStatusProvider(ctx context.Context, clients []*rpc.Client, bakers []string, statusChannel chan<- common.ProviderStatusUpdatedReport) {
	blockChannel := make(chan *rpc.BlockHeaderLogEntry)
	id, err := uuid.NewRandom()
	if err != nil {
		slog.Error("failed to generate block subscriber (baker status provider) uuid", "error", err.Error())
		return
	}
	blockSubscribers[id] = blockChannel

	go func() {
		defer func() {
			delete(blockSubscribers, id)
			close(blockChannel)
		}()

		if blockProviderCount.Load() == 0 {
			slog.Warn("no block providers are running, rights provider will not work until at least one block provider is running")
		}

		lastProcessedBlockHeight = int64(0)

		status := BakersStatus{
			Level:  0,
			Bakers: map[string]*rpc.Delegate{},
		}

		for {
			select {
			case <-ctx.Done():
				return
			case block, ok := <-blockChannel:
				if !ok {
					// levelChannel is closed, exit the loop
					return
				}

				if ctx.Done() != nil {
					return
				}

				if status.Level >= block.Level {
					continue
				}

				syncedClients := lo.Filter(clients, func(client *rpc.Client, _ int) bool {
					status, err := client.GetStatus(ctx)
					return status.SyncState == "synced" || (err != nil && strings.Contains(err.Error(), "status 403"))
				})

				for _, baker := range bakers {
					status.Bakers[baker], _ = getBakerStatusFor(ctx, syncedClients, baker)
				}
				status.Level = block.Level
				statusChannel <- &BakersStatusUpdate{
					Status: status,
				}
			}
		}
	}()
}
