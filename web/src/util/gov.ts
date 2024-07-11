import { BLOCK_TIME, PROPOSAL_PERIOD_QUORUM_PERCENTAGE, PROPOSAL_REQUIRED_MAJORITY_PERCENTAGE } from "@src/common/constants";
import type { NodeStatus, VotingPeriodInfo } from "@src/common/types/status";
import type { ExtendGovernanceExplorationOrPromotionPeriodDetail, ExtendedGovernanceProposalPeriodDetail, GovernanceExplorationOrPromotionPeriodDetail, GovernanceProposalPeriodDetail } from "@src/common/types/governance";
import { formatDistanceToNow } from "date-fns/formatDistanceToNow";

export function getVotingPeriodEndDate(info: VotingPeriodInfo) {
	if (!info) {
		return new Date(0);
	}

	const remainingMs = (info?.remaining ?? 1) * BLOCK_TIME;
	const endDate = new Date(Date.now() + remainingMs);
	return endDate;
}

export function getVotingPeriodTimeLeft(info?: VotingPeriodInfo) {
	if (!info) {
		return "N/A";
	}
	const endDate = getVotingPeriodEndDate(info);
	return formatDistanceToNow(endDate);
}

export function getGovernancePeriodBlocksLeft(info: VotingPeriodInfo) {
	if (!info) {
		return "N/A";
	}
	return info.remaining;
}


export function pickVotingPeriodInfo(statuses: Array<NodeStatus>) {
	for (const status of statuses) {
		if (status.block?.voting_period_info) {
			return status.block.voting_period_info
		}
	}
	return undefined
}

export function getCurrentBlock(statuses: Array<NodeStatus>) {
	for (const status of statuses) {
		if (status.block) {
			return status.block.level_info.level
		}
	}
	return undefined
}

export function formatVotes(votes: string) {
	return (BigInt(votes) / (10n ** 6n)).toLocaleString()
}


export function preprocessProposalPeriodDetailProposals(proposalPeriod: GovernanceProposalPeriodDetail): ExtendedGovernanceProposalPeriodDetail {
	const { proposals, voters } = proposalPeriod

	const totalVotes = voters.reduce((acc, voter) => {
		return acc + BigInt(voter.voting_power)
	}, 0n)

	return {
		...proposalPeriod,
		proposals: proposals.map((proposal) => {
			const percentage = Number(BigInt(proposal.upvotes) * 100000n / totalVotes) / 1000

			return {
				...proposal,
				upvotes: formatVotes(proposal.upvotes),
				percentage: Math.round(percentage * 100) / 100,
				percentageString: `${percentage}%`,
				quorumReached: BigInt(proposal.upvotes) > (totalVotes / 100n * PROPOSAL_PERIOD_QUORUM_PERCENTAGE)
			}
		}).sort((a, b) => (BigInt(b.upvotes) > BigInt(a.upvotes) ? 1 : -1))
	}
}

export function preprocessExpPromPeriodDetail(period: GovernanceExplorationOrPromotionPeriodDetail): ExtendGovernanceExplorationOrPromotionPeriodDetail {
	const { voters } = period

	const totalVotes = voters.reduce((acc, voter) => {
		return acc + BigInt(voter.voting_power)
	}, 0n)

	const totalYayNay = BigInt(period.summary.yay) + BigInt(period.summary.nay)
	const totalCasted = totalYayNay + BigInt(period.summary.pass)

	return {
		...period,
		summary: {
			...period.summary,
			yayPercentage: Math.round(Number(BigInt(period.summary.yay) * 100000n / totalVotes) / 10) / 100,
			nayPercentage: Math.round(Number(BigInt(period.summary.nay) * 100000n / totalVotes) / 10) / 100,
			passPercentage: Math.round(Number(BigInt(period.summary.pass) * 100000n / totalVotes) / 10) / 100,
		},
		quorumReached: totalCasted > (totalVotes * BigInt(period.quorum) / 10000n),
		quorumPercentage: Math.min(Number(BigInt(period.summary.yay + period.summary.nay + period.summary.pass) * 100000n / totalVotes) / 1000, 100),
		requiredQuorumPercentage: Number(period.quorum) / 100,
		majorityReached: BigInt(period.summary.yay) >= (totalYayNay * PROPOSAL_REQUIRED_MAJORITY_PERCENTAGE) / 100n,
		majorityPercentage: Math.round(Number(BigInt(period.summary.yay) * 100000n / totalYayNay) / 10) / 100,
		requiredMajorityPercentage: Number(PROPOSAL_REQUIRED_MAJORITY_PERCENTAGE),
	}

}