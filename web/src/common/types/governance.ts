import type { VotingPeriodInfo } from "./status"

export type Voter = {
	pkh: string,
	voting_power: number
}

export type BallotVote = "yay" | "nay" | "pass";

export type Ballot = {
	pkh: string,
	ballot: BallotVote
}

export type GovernanceProposalPeriodDetailProposal = {
	proposal: string, upvotes: string
}

export type GovernanceProposalPeriodDetail = {
	info: VotingPeriodInfo & {
		voting_period: {
			kind: "proposal"
		}
	}
	voters: Array<Voter>
	proposals: Array<GovernanceProposalPeriodDetailProposal>
	votes: { [key: string]: string[] }
}

export type ExtendedGovernanceProposalPeriodDetail = GovernanceProposalPeriodDetail & {
	proposals: Array<GovernanceProposalPeriodDetailProposal & {
		upvotes: string,
		percentage: number,
		percentageString: string,
		quorumReached: boolean
	}>
}

export type GovernanceExplorationOrPromotionPeriodDetail = {
	info: VotingPeriodInfo & {
		voting_period: {
			kind: "exploration" | "promotion"
		}
	}
	voters: Array<Voter>
	summary: {
		yay: string,
		nay: string,
		pass: string,
	}
	proposal: string
	// Returned value is percent * 10000 i.e. 5820 for 58.20%.
	quorum: number
	ballots: Array<Ballot>
}

export type CommonPeriodDetail = {
	info: VotingPeriodInfo & {
		voting_period: {
			kind: "testing" | "adoption"
		}
	},
	voters: Array<Voter>
	proposal: string
}

export const testExpDetail: GovernanceExplorationOrPromotionPeriodDetail = {
	info: {
		position: 100,
		remaining: 100,
		voting_period: {
			index: 10,
			kind: "exploration",
			start_position: 50
		}
	},
	voters: [
		{ pkh: "tz1P6WKJu2rcbxKiKRZHKQKmKrpC9TfW1AwM", voting_power: 400 },
		{ pkh: "tz1hZvgjekGo7DmQjWh7XnY5eLQD8wNYPczE", voting_power: 200 },
	],
	summary: {
		yay: "110",
		nay: "300",
		pass: "200"
	},
	proposal: "Pt1JoinAscentToMountVinsonAGNqxgMLDAB8TqZpDwMTU5eCx",
	quorum: 5011,
	ballots: [
		{ pkh: "tz1P6WKJu2rcbxKiKRZHKQKmKrpC9TfW1AwM", ballot: "yay" },
		{ pkh: "tz1hZvgjekGo7DmQjWh7XnY5eLQD8wNYPczE", ballot: "nay" },
	]
}

export type ExtendGovernanceExplorationOrPromotionPeriodDetail = GovernanceExplorationOrPromotionPeriodDetail & {
	summary: {
		yay: string,
		yayPercentage: number,
		nay: string,
		nayPercentage: number,
		pass: string,
		passPercentage: number,
	}
	quorumReached: boolean
	quorumPercentage: number
	requiredQuorumPercentage: number
	majorityReached: boolean
	majorityPercentage: number
	requiredMajorityPercentage: number
}


export type GovernancePeriodDetail = GovernanceProposalPeriodDetail | GovernanceExplorationOrPromotionPeriodDetail | CommonPeriodDetail