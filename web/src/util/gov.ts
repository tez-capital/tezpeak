import type { NodeStatus, NormalizedVotingPeriodInfo } from "@src/common/types";
// import { formatDistanceStrict, formatDistance } from "date-fns"

export const pickVotingPeriodInfo = (statuses: Array<NodeStatus>) => {
	for (const status of statuses) {
		if (status.block?.voting_period_info) {
			return {
				kind: status.block.voting_period_info.voting_period.kind,
				index: status.block.voting_period_info.voting_period.index,
				remaining: status.block.voting_period_info.remaining,
			} as NormalizedVotingPeriodInfo
		}
	}
	return undefined
}