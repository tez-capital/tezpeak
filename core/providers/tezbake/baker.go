package tezbake

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"

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

func (s *BakersStatusUpdate) GetData() any {
	return s.Status
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

func getDelegateDelegatedContracts(ctx context.Context, client *common.ActiveRpcNode, addr tezos.Address, id rpc.BlockID) ([]tezos.Address, error) {
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

func getDelegateStakingStatusFromRawContext(ctx context.Context, client *common.ActiveRpcNode, delegate tezos.Address, id rpc.BlockID) (*BakerStakingStatus, error) {
	u := fmt.Sprintf("chains/main/blocks/%s/context/raw/json/staking_balance/%s", id, delegate)

	var status BakerStakingStatus
	err := client.Get(ctx, u, &status)
	return &status, err
}

func getBakerStatusFor(ctx context.Context, baker string) (*BakerStakingStatus, error) {
	addr, err := tezos.ParseAddress(baker)
	if err != nil {
		return nil, err
	}
	status, err := common.AttemptWithRpcClients(ctx, func(client *common.ActiveRpcNode) (*BakerStakingStatus, error) {
		if !client.IsGovernanceProvider {
			return nil, errors.New("node is not a rights provider")
		}
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

func setupBakerStatusProviders(ctx context.Context, bakers []string, statusChannel chan<- common.StatusUpdatedReport) {
	blockChannelId, blockChannel, err := common.SubscribeToBlockHeaderEvents()
	if err != nil {
		slog.Error("failed to subscribe to block events", "error", err.Error())
		return
	}

	go func() {
		defer func() {
			common.UnsubscribeFromBlockHeaderEvents(blockChannelId)
		}()

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

				for _, baker := range bakers {
					status.Bakers[baker], _ = getBakerStatusFor(ctx, baker)
				}
				status.Level = block.Level
				statusChannel <- &BakersStatusUpdate{
					Status: status,
				}
			}
		}
	}()
}
