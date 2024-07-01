package providers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"sync/atomic"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/tez-capital/tezpeak/core/common"
	"github.com/trilitech/tzgo/rpc"
	"github.com/trilitech/tzgo/tezos"
)

var (
	rightsProviderCount atomic.Int64
	rightsChannel       = make(chan *RightsStatus, 100)
)

type BlockRights struct {
	Level           int64            `json:"level"`
	Rights          map[string][]int `json:"rights"`
	RealizedChecked bool             `json:"realized_checked"`
}

type RightsStatus struct {
	Level  int64          `json:"level"`
	Rights []*BlockRights `json:"rights"`
}

type RightsStatusUpdate struct {
	Status RightsStatus
}

func (s *RightsStatusUpdate) GetId() string {
	return "rights"
}

func (s *RightsStatusUpdate) GetData() interface{} {
	return s.Status
}

func (s *RightsStatusUpdate) GetKind() common.StatusUpdateKind {
	return common.RightsStatusUpdateKind
}

func initRights(bakers []string) (map[string]int, map[string]int) {
	attestations := map[string]int{}
	baking := map[string]int{}
	for _, baker := range bakers {
		attestations[baker] = 0
		baking[baker] = 0
	}
	return baking, attestations
}

// [{"level":5026842,"delegates":[{"delegate":"tz1P6WKJu2rcbxKiKRZHKQKmKrpC9TfW1AwM","first_slot":2608,"attestation_power":1,"consensus_key":"tz1P6WKJu2rcbxKiKRZHKQKmKrpC9TfW1AwM"}]}]
type rights struct {
	Level     int64 `json:"level"`
	Delegates []struct {
		Delegate         string `json:"delegate"`
		FirstSlot        int64  `json:"first_slot"`
		AttestationPower int    `json:"attestation_power"`
		ConsensusKey     string `json:"consensus_key"`
	} `json:"delegates"`
}

func getBlockRights(ctx context.Context, clients []*rpc.Client, block int64) (rights, rights, error) {
	bakingRights := rights{}
	attestationRights := rights{}
	var bakingRightsErr, attestationRightsErr error
	bakingRightsChan := make(chan struct{})
	attestationRightsChan := make(chan struct{})

	go func() {
		url := fmt.Sprintf("chains/main/blocks/head/helpers/baking_rights?all=true&max_priority=1&level=%d", block)
		bakingRights, bakingRightsErr = attemptWithClients(clients, func(client *rpc.Client) (rights, error) {
			bakingRights := make([]rights, 0)
			err := client.Get(ctx, url, &bakingRights)
			result := rights{Level: block}
			if len(bakingRights) > 0 {
				result = bakingRights[0]
			}
			return result, err
		})
		bakingRightsChan <- struct{}{}
	}()

	go func() {
		url := fmt.Sprintf("chains/main/blocks/head/helpers/attestation_rights?all=true&max_priority=1&level=%d", block)
		attestationRights, attestationRightsErr = attemptWithClients(clients, func(client *rpc.Client) (rights, error) {
			attestationRights := make([]rights, 0)
			err := client.Get(ctx, url, &attestationRights)
			result := rights{Level: block}
			if len(attestationRights) > 0 {
				result = attestationRights[0]
			}
			return result, err
		})
		attestationRightsChan <- struct{}{}
	}()
	<-bakingRightsChan
	<-attestationRightsChan
	return bakingRights, attestationRights, errors.Join(attestationRightsErr, bakingRightsErr)
}

func getBlockRightsFor(ctx context.Context, clients []*rpc.Client, block int64, bakers []string) (*BlockRights, error) {
	relevantBakingRights, relevantAttestationRights := initRights(bakers)

	bakingRights, attestationRights, err := getBlockRights(ctx, clients, block-1)

	for _, right := range bakingRights.Delegates {
		if _, ok := relevantBakingRights[right.Delegate]; !ok {
			continue
		}
		relevantBakingRights[right.Delegate]++
	}

	for _, right := range attestationRights.Delegates {
		if _, ok := relevantAttestationRights[right.Delegate]; !ok {
			continue
		}
		relevantAttestationRights[right.Delegate] += right.AttestationPower
	}

	if err != nil {
		slog.Warn("Reported error while getting block rights", "error", err.Error())
	}

	rights := map[string][]int{}
	for _, baker := range bakers {
		rights[baker] = []int{relevantBakingRights[baker], relevantAttestationRights[baker]}
	}

	return &BlockRights{
		Level:  block,
		Rights: rights,
	}, nil
}

func checkRealized(ctx context.Context, clients []*rpc.Client, rights *BlockRights) (*BlockRights, error) {
	if rights.RealizedChecked {
		return rights, nil
	}

	hasAnyRights := false
	for _, r := range rights.Rights {
		if len(r) > 1 && (r[0] > 0 || r[1] > 0) {
			hasAnyRights = true
			break
		}
	}
	if !hasAnyRights {
		rights.RealizedChecked = true
		return rights, nil
	}

	header, err := attemptWithClients(clients, func(client *rpc.Client) (*rpc.BlockHeader, error) {
		return client.GetBlockHeader(ctx, rpc.BlockLevel(rights.Level))
	})

	if err != nil {
		return rights, err
	}

	ops, err := attemptWithClients(clients, func(client *rpc.Client) ([][]rpc.Operation, error) {
		ops, err := client.GetBlockOperations(ctx, rpc.BlockLevel(rights.Level))
		return ops, err
	})
	if err != nil {
		return rights, err
	}

	validAttestations := lo.Reduce(ops, func(acc []string, g []rpc.Operation, _ int) []string {
		for _, tx := range g {
			for _, c := range tx.Contents {
				if c.Kind() == tezos.OpTypeAttestation {
					acc = append(acc, c.Meta().Delegate.String())
				}
			}
		}
		return acc
	}, []string{})

	for baker, r := range rights.Rights {
		if len(r) > 1 {
			blockRights, attestationRights := r[0], r[1]
			bakedBlock := 0
			if blockRights > 0 && header.PayloadRound == 0 {
				bakedBlock = 1
			}

			attestedBlock := 0
			for _, attester := range validAttestations {
				if attester == baker {
					attestedBlock = 1
					break
				}
			}

			rights.Rights[baker] = []int{blockRights, attestationRights, bakedBlock, attestedBlock}
		}
	}
	rights.RealizedChecked = true

	return rights, nil
}

func StartRightsStatusProvider(ctx context.Context, clients []*rpc.Client, bakers []string, window int64, statusChannel chan<- common.ProviderStatusUpdatedReport) {
	blockChannel := make(chan *rpc.BlockHeaderLogEntry)
	id, err := uuid.NewRandom()
	if err != nil {
		slog.Error("failed to generate block subscriber (rights status provider) uuid", "error", err.Error())
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

		status := RightsStatus{
			Level:  0,
			Rights: []*BlockRights{},
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

				// get slice of levels to query
				minLevel := max(0, block.Level-window/2)
				maxLevel := block.Level + window/2
				newRights := []*BlockRights{}
				lastCachedLevel := int64(0)
				for _, right := range status.Rights {
					if right.Level < minLevel || right.Level > maxLevel {
						continue
					}
					newRights = append(newRights, right)
					lastCachedLevel = right.Level
				}

				for i := max(lastCachedLevel+1, minLevel); i < maxLevel; i++ {
					rights, err := getBlockRightsFor(ctx, syncedClients, i, bakers)
					if err != nil {
						slog.Error("failed to get block rights", "error", err.Error())
						continue
					}
					newRights = append(newRights, rights)
				}

				for _, right := range newRights {
					if right.Level > block.Level { // we do not check future rights
						break
					}
					checkRealized(ctx, syncedClients, right)
				}

				status.Level = block.Level
				status.Rights = newRights

				rightsChannel <- &status
				statusChannel <- &RightsStatusUpdate{
					Status: status,
				}
			}
		}
	}()
}
