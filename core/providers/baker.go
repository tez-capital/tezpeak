package providers

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/tez-capital/tezpeak/constants"
	"github.com/tez-capital/tezpeak/core/common"
	"github.com/trilitech/tzgo/rpc"
	"github.com/trilitech/tzgo/tezos"
)

type BakersStatus struct {
	Level  int64                          `json:"level"`
	Bakers map[string]*BakerStakingStatus `json:"bakers"`
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

type BakerStakingStatus struct {
	Balance                  string  `json:"balance"`
	OwnFrozen                tezos.Z `json:"own_frozen"`
	StakedBalance            string  `json:"staked_balance"`
	StakedFrozen             tezos.Z `json:"staked_frozen"`
	ExternalStakedBalance    string  `json:"external_staked_balance"`
	Delegated                tezos.Z `json:"delegated"`
	ExternalDelegatedBalance string  `json:"external_delegated_balance"`
	DelegatorsCount          int     `json:"delegators_count"`
}

func getDelegateDelegatedContracts(ctx context.Context, client *rpc.Client, addr tezos.Address, id rpc.BlockID) ([]tezos.Address, error) {
	u := fmt.Sprintf("chains/main/blocks/%s/context/delegates/%s/delegated_contracts", id, addr)

	var delegatedContracts []tezos.Address
	err := client.Get(ctx, u, &delegatedContracts)
	if err != nil {
		if strings.Contains(err.Error(), "delegate.not_registered") {
			return []tezos.Address{}, constants.ErrDelegateNotRegistered
		}
		var rpcErrors []rpc.GenericError
		err2 := client.Get(ctx, u, &rpcErrors)
		if err2 != nil {
			return nil, err
		}
		for _, rpcError := range rpcErrors {
			if strings.Contains(rpcError.ID, "delegate.not_registered") {
				return []tezos.Address{}, constants.ErrDelegateNotRegistered
			}
		}
	}

	return delegatedContracts, err
}

func getDelegateStakingStatusFromRawContext(ctx context.Context, client *rpc.Client, delegate tezos.Address, id rpc.BlockID) (*BakerStakingStatus, error) {
	u := fmt.Sprintf("chains/main/blocks/%s/context/raw/json/staking_balance/%s", id, delegate)

	var status BakerStakingStatus
	err := client.Get(ctx, u, &status)
	return &status, err
}

func getBakerStatusFor(ctx context.Context, clients []*rpc.Client, baker string) (*BakerStakingStatus, error) {
	addr, err := tezos.ParseAddress(baker)
	if err != nil {
		return nil, err
	}
	status, err := attemptWithClients(clients, func(client *rpc.Client) (*BakerStakingStatus, error) {
		acc, err := getDelegateStakingStatusFromRawContext(ctx, client, addr, rpc.Head)
		if err != nil {
			return nil, err
		}
		balance, err := client.GetContractBalance(ctx, addr, rpc.Head)
		if err != nil {
			return nil, err
		}
		acc.Balance = balance.String()
		delegators, err := getDelegateDelegatedContracts(ctx, client, addr, rpc.Head)
		if err != nil {
			return nil, err
		}
		acc.DelegatorsCount = len(delegators)
		acc.StakedBalance = acc.OwnFrozen.String()
		acc.ExternalStakedBalance = acc.StakedFrozen.String()
		acc.ExternalDelegatedBalance = acc.Delegated.Sub(balance).String()

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
			Bakers: map[string]*BakerStakingStatus{},
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
