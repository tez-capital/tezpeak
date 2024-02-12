package providers

import (
	"context"
	"log/slog"

	"blockwatch.cc/tzgo/rpc"
	"github.com/google/uuid"
	"github.com/tez-capital/tezpeak/core/common"
)

type LedgerStatus struct {
	Level int64 `json:"level"`
}

func StartLedgerStatusProvider(ctx context.Context, clients []*rpc.Client, bakers []string, window int64, statusChannel chan<- common.ProviderStatusUpdatedReport) {
	rightsChannel := make(chan *RightsStatus)
	id, err := uuid.NewRandom()
	if err != nil {
		slog.Error("failed to generate block subscriber (rights status provider) uuid", "error", err.Error())
		return
	}
	rightsSubscribers[id] = rightsChannel

	go func() {
		defer func() {
			delete(rightsSubscribers, id)
			close(rightsChannel)
		}()

		if rightsProviderCount.Load() == 0 {
			slog.Warn("no block providers are running, rights provider will not work until at least one block provider is running")
		}

		status := LedgerStatus{
			Level: 0,
		}

		for {
			select {
			case <-ctx.Done():
				return
			case rightsStatus, ok := <-rightsChannel:
				if !ok {
					// levelChannel is closed, exit the loop
					return
				}

				if ctx.Done() != nil {
					return
				}

				if status.Level >= rightsStatus.Level {
					continue
				}

				status.Level = rightsStatus.Level

			}
		}

		// for {
		// 	select {
		// 	case <-ctx.Done():
		// 		return
		// 	case block, ok := <-blockChannel:
		// 		if !ok {
		// 			// levelChannel is closed, exit the loop
		// 			return
		// 		}

		// 		if ctx.Done() != nil {
		// 			return
		// 		}

		// 		if status.Level >= block.Level {
		// 			continue
		// 		}

		// 		syncedClients := lo.Filter(clients, func(client *rpc.Client, _ int) bool {
		// 			status, err := client.GetStatus(ctx)
		// 			return status.SyncState == "synced" || (err != nil && strings.Contains(err.Error(), "status 403"))
		// 		})

		// 		// get slice of levels to query
		// 		minLevel := max(0, block.Level-window/2)
		// 		maxLevel := block.Level + window/2
		// 		newRights := []BlockRights{}
		// 		lastCachedLevel := int64(0)
		// 		for _, right := range status.Rights {
		// 			if right.Level < minLevel || right.Level > maxLevel {
		// 				continue
		// 			}
		// 			newRights = append(newRights, right)
		// 			lastCachedLevel = right.Level
		// 		}

		// 		for i := max(lastCachedLevel+1, minLevel); i < maxLevel; i++ {
		// 			rights, err := getBlockRightsFor(ctx, syncedClients, i, bakers)
		// 			if err != nil {
		// 				slog.Error("failed to get block rights", "error", err.Error())
		// 				continue
		// 			}
		// 			newRights = append(newRights, *rights)
		// 		}

		// 		status.Level = block.Level
		// 		status.Rights = newRights

		// 		for _, right := range status.Rights {
		// 			if right.Level > block.Level { // we do not check future rights
		// 				break
		// 			}
		// 			checkRealized(ctx, syncedClients, &right)
		// 		}
		// 		rightsChannel <- &status
		// 		statusChannel <- &RightsStatusUpdate{
		// 			Status: status,
		// 		}
		// 	}
		// }
	}()
}
