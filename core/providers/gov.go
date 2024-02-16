package providers

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"blockwatch.cc/tzgo/codec"
	"blockwatch.cc/tzgo/rpc"
	"blockwatch.cc/tzgo/signer/remote"
	"blockwatch.cc/tzgo/tezos"
	"github.com/tez-capital/tezpeak/configuration"
	"github.com/tez-capital/tezpeak/constants"
	"github.com/tez-capital/tezpeak/util"
)

type GovernanceProvider struct {
	configuration *configuration.Runtime
	nodes         []*rpc.Client
	signer        string
}

type GovernancePeriodDetail struct {
	Info      *rpc.VotingPeriodInfo `json:"info"`
	Voters    rpc.VoterList         `json:"voters"`
	Summary   *rpc.BallotSummary    `json:"summary"`
	Proposal  tezos.ProtocolHash    `json:"proposal"`
	Proposals rpc.ProposalList      `json:"proposals"`
	Quorum    int                   `json:"quorum"`
	Ballots   rpc.BallotList        `json:"ballots"`
}

type UpvoteParams struct {
	Source    tezos.Address        `json:"source"`
	Proposals []tezos.ProtocolHash `json:"proposals"`
	Period    int32                `json:"period"`
}

func (p *UpvoteParams) ToContents() codec.Operation {
	return &codec.Proposals{
		Source:    p.Source,
		Proposals: p.Proposals,
		Period:    p.Period,
	}
}

type VoteParams struct {
	Source   tezos.Address      `json:"source"`
	Proposal tezos.ProtocolHash `json:"proposal"`
	Period   int32              `json:"period"`
	Ballot   string             `json:"ballot"`
}

func (p *VoteParams) ToContents() codec.Operation {
	return &codec.Ballot{
		Source:   p.Source,
		Proposal: p.Proposal,
		Period:   p.Period,
		Ballot:   tezos.ParseBallotVote(p.Ballot),
	}
}

func InitGovernanceProvider(ctx context.Context, configuration *configuration.Runtime) *GovernanceProvider {
	bakerNodeClient, _ := rpc.NewClient(configuration.Node, nil)

	rightProviderRpcs := make([]*rpc.Client, 0, len(configuration.ReferenceNodes)+1)
	rightProviderRpcs = append(rightProviderRpcs, bakerNodeClient)

	for id, node := range configuration.ReferenceNodes {
		if node.Address == "" {
			slog.Warn("no address for node", "id", id)
			continue
		}

		if node.IsGovernanceProvider {
			if client, err := rpc.NewClient(node.Address, nil); err == nil {
				rightProviderRpcs = append(rightProviderRpcs, client)
			} else {
				slog.Debug("failed to connect to node", "source", node.Address, "error", err.Error())
			}
		}
	}

	return &GovernanceProvider{
		nodes:         rightProviderRpcs,
		configuration: configuration,
		signer:        configuration.Signer,
	}
}

func (governanceProvider *GovernanceProvider) CanVote() bool {
	return governanceProvider.configuration.Mode == configuration.PrivatePeakMode
}

func wrapInWaithGroup(wg *sync.WaitGroup, f func()) {
	wg.Add(1)
	go func() {
		f()
		wg.Done()
	}()
}

func (governanceProvider *GovernanceProvider) startVotersCollector(ctx context.Context, detail *GovernancePeriodDetail, wg *sync.WaitGroup) {
	wrapInWaithGroup(wg, func() {
		voters, _ := attemptWithClients(governanceProvider.nodes, func(client *rpc.Client) (rpc.VoterList, error) {
			voters, err := client.ListVoters(ctx, rpc.Head)
			if err != nil {
				return nil, err
			}
			// voters[0].Power <- power
			return voters, err
		})

		detail.Voters = voters
	})
}

func (governanceProvider *GovernanceProvider) startProposalsCollector(ctx context.Context, detail *GovernancePeriodDetail, wg *sync.WaitGroup) {
	wrapInWaithGroup(wg, func() {
		proposals, _ := attemptWithClients(governanceProvider.nodes, func(client *rpc.Client) (rpc.ProposalList, error) {
			proposals, err := client.ListProposals(ctx, rpc.Head)
			if err != nil {
				return nil, err
			}
			return proposals, err
		})
		detail.Proposals = proposals
	})
}

func (governanceProvider *GovernanceProvider) startQuorumCollector(ctx context.Context, detail *GovernancePeriodDetail, wg *sync.WaitGroup) {
	wrapInWaithGroup(wg, func() {
		quorum, _ := attemptWithClients(governanceProvider.nodes, func(client *rpc.Client) (int, error) {
			quorum, err := client.GetVoteQuorum(ctx, rpc.Head)
			if err != nil {
				return 0, err
			}
			return quorum, err
		})
		detail.Quorum = quorum
	})
}

func (governanceProvider *GovernanceProvider) startBallotsCollector(ctx context.Context, detail *GovernancePeriodDetail, wg *sync.WaitGroup) {
	wrapInWaithGroup(wg, func() {
		ballotList, _ := attemptWithClients(governanceProvider.nodes, func(client *rpc.Client) (rpc.BallotList, error) {
			ballots, err := client.ListBallots(ctx, rpc.Head)
			if err != nil {
				return nil, err
			}
			return ballots, err
		})
		detail.Ballots = ballotList
	})
}

func (governanceProvider *GovernanceProvider) startSummaryCollector(ctx context.Context, detail *GovernancePeriodDetail, wg *sync.WaitGroup) {
	wrapInWaithGroup(wg, func() {
		summary, _ := attemptWithClients(governanceProvider.nodes, func(client *rpc.Client) (*rpc.BallotSummary, error) {
			summary, err := client.GetVoteResult(ctx, rpc.Head)
			if err != nil {
				return nil, err
			}
			return &summary, err
		})
		detail.Summary = summary
	})
}

func (governanceProvider *GovernanceProvider) startProtocolCollector(ctx context.Context, detail *GovernancePeriodDetail, wg *sync.WaitGroup) {
	wrapInWaithGroup(wg, func() {
		currentProposal, _ := attemptWithClients(governanceProvider.nodes, func(client *rpc.Client) (tezos.ProtocolHash, error) {
			proposal, err := client.GetVoteProposal(ctx, rpc.Head)
			if err != nil {
				return tezos.ProtocolHash{}, err
			}
			return proposal, err
		})
		detail.Proposal = currentProposal
	})
}

func (governanceProvider *GovernanceProvider) GetAvailablePkhs(ctx context.Context) ([]string, error) {
	return governanceProvider.configuration.Bakers, nil
}

func (governanceProvider *GovernanceProvider) GetGovernancePeriodDetail(ctx context.Context) (*GovernancePeriodDetail, error) {
	periodInfo, err := attemptWithClients(governanceProvider.nodes, func(client *rpc.Client) (*rpc.VotingPeriodInfo, error) {
		meta, err := client.GetBlockMetadata(ctx, rpc.Head)
		if err != nil {
			return nil, err
		}
		return meta.VotingPeriodInfo, err
	})
	if err != nil {
		slog.Warn("failed to get voting period kind", "error", err.Error())
		return nil, err
	}

	detail := &GovernancePeriodDetail{
		Info: periodInfo,
	}

	var wg sync.WaitGroup

	if periodInfo.VotingPeriod.Kind == tezos.VotingPeriodProposal ||
		periodInfo.VotingPeriod.Kind == tezos.VotingPeriodExploration ||
		periodInfo.VotingPeriod.Kind == tezos.VotingPeriodPromotion {
		governanceProvider.startVotersCollector(ctx, detail, &wg)
	}

	if periodInfo.VotingPeriod.Kind == tezos.VotingPeriodProposal {
		governanceProvider.startProposalsCollector(ctx, detail, &wg)
	}

	if periodInfo.VotingPeriod.Kind == tezos.VotingPeriodExploration || periodInfo.VotingPeriod.Kind == tezos.VotingPeriodPromotion {
		governanceProvider.startQuorumCollector(ctx, detail, &wg)
		governanceProvider.startBallotsCollector(ctx, detail, &wg)
		governanceProvider.startSummaryCollector(ctx, detail, &wg)
		governanceProvider.startProtocolCollector(ctx, detail, &wg)
	}

	wg.Wait()
	return detail, nil
}

func (governanceProvider *GovernanceProvider) buildAndBroadcastGovernanceOperation(ctx context.Context, pkh tezos.Address, contents codec.Operation) (tezos.OpHash, error) {
	rs, err := remote.New(governanceProvider.signer, nil)
	if err != nil {
		err = util.TryUnwrapRPCError(err)
		slog.Error("failed to create remote signer", "error", err.Error())
		return tezos.OpHash{}, errors.Join(constants.ErrFailedToCreateRemoteSigner, err)
	}

	key, err := rs.GetKey(ctx, pkh)
	if err != nil {
		err = util.TryUnwrapRPCError(err)
		slog.Error("failed to get key", "error", err.Error())
		return tezos.OpHash{}, errors.Join(constants.ErrFailedToGetPublicKey, err)
	}

	// complete the operation
	op, err := attemptWithClients(governanceProvider.nodes, func(client *rpc.Client) (*codec.Op, error) {
		params, err := client.GetParams(ctx, rpc.Head)
		if err != nil {
			return nil, err
		}
		op := codec.NewOp().WithContents(contents).WithSource(pkh)
		op.WithTTL(constants.MAX_OPERATION_TTL)
		op.WithContents(contents)

		op = op.WithParams(params)
		err = client.Complete(ctx, op, key)
		if err != nil {
			return op, err
		}

		return op, err
	})
	if err != nil {
		err = util.TryUnwrapRPCError(err)
		slog.Error("failed to complete operation", "error", err.Error())
		return tezos.OpHash{}, errors.Join(constants.ErrFailedToCompleteOperation, err)
	}

	signature, err := rs.SignOperation(ctx, pkh, op)
	if err != nil {
		err = util.TryUnwrapRPCError(err)
		slog.Error("failed to sign operation", "error", err.Error())
		return tezos.OpHash{}, errors.Join(constants.ErrFailedToSignOperation, err)
	}
	op = op.WithSignature(signature)

	opHash, err := attemptWithClients(governanceProvider.nodes, func(client *rpc.Client) (tezos.OpHash, error) {
		opHash, err := client.Broadcast(ctx, op)
		if err != nil {
			slog.Error("failed to broadcast operation", "error", err.Error())
			return tezos.OpHash{}, err
		}
		return opHash, err
	})
	if err != nil {
		err = util.TryUnwrapRPCError(err)
		slog.Error("failed to broadcast operation", "error", err.Error())
		return tezos.OpHash{}, errors.Join(constants.ErrFailedToBroadcastOperation, err)
	}

	return opHash, nil
}

func (governanceProvider *GovernanceProvider) Upvote(ctx context.Context, params *UpvoteParams) (tezos.OpHash, error) {
	return governanceProvider.buildAndBroadcastGovernanceOperation(ctx, params.Source, params.ToContents())
}

func (governanceProvider *GovernanceProvider) Vote(ctx context.Context, params *VoteParams) (tezos.OpHash, error) {
	return governanceProvider.buildAndBroadcastGovernanceOperation(ctx, params.Source, params.ToContents())
}

func (governanceProvider *GovernanceProvider) WaitConfirmation(ctx context.Context, opHash string) (bool, error) {
	op, err := tezos.ParseOpHash(opHash)
	if err != nil {
		return false, err
	}

	ctx, cancel := context.WithTimeout(ctx, constants.MAX_WAIT_FOR_CONFIRMATION*time.Second)
	defer cancel()
	result, err := attemptWithClients(governanceProvider.nodes, func(client *rpc.Client) (bool, error) {
		result := rpc.NewResult(op)
		client.Listen()
		defer client.Close()
		result.Listen(client.BlockObserver)
		result.WaitContext(ctx)

		if err := result.Err(); err != nil {
			return false, err
		}

		return true, err
	})
	return result, err
}
