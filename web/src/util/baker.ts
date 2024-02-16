import { DEPOSIT_LIMIT_STAKING_CAPACITY_MULTIPLIER } from "@src/common/constants";
import type { BakerStatus, BlockRights, NormalizedBlockRights } from "@src/common/types/status";

export function calculateFreeSpace(baker: BakerStatus) {
	let limit = BigInt(baker.frozen_deposits_limit)
	if (limit === 0n) {
		limit = BigInt(baker.full_balance)
	}

	const delegated = BigInt(baker.delegated_balance)

	return (limit * DEPOSIT_LIMIT_STAKING_CAPACITY_MULTIPLIER) - delegated
}

export function calculateStakingCapacity(baker: BakerStatus) {
	let limit = BigInt(baker.frozen_deposits_limit)
	if (limit === 0n) {
		limit = BigInt(baker.full_balance)
	}

	return limit * DEPOSIT_LIMIT_STAKING_CAPACITY_MULTIPLIER
}

export function getBakerColor(str: string) {
	let hash = 0;
	str.split('').forEach(char => {
		hash = char.charCodeAt(0) + ((hash << 5) - hash)
	})
	let color = '#'
	for (let i = 0; i < 3; i++) {
		const value = (hash >> (i * 8)) & 0xff
		color += value.toString(16).padStart(2, '0')
	}
	return color
}

export function normalizeBlockRights(blockRights: BlockRights) {
	const result = [] as Array<NormalizedBlockRights>
	for (const [baker, rights] of Object.entries(blockRights.rights)) {
		result.push({
			baker,
			blocks: rights[0],
			attestations: rights[1],
			realizedBlocks: rights.length > 2 ? rights[2] : 0,
			realizedAttestations: rights.length > 3 ? rights[3] : 0,
		})
	}
	return result
}