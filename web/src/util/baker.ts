import type { BlockRights } from "@src/common/types/status";

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

// done/missed only once the backend checked the level (realized_checked), otherwise pending
export type RightStatus = "upcoming" | "pending" | "done" | "missed"

export type BakerLevelRights = {
	level: number
	blocks: number
	attestations: number
	blockStatus?: RightStatus
	attestationStatus?: RightStatus
}

function rightStatus(assigned: number, future: boolean, checked: boolean, realized?: number): RightStatus | undefined {
	if (!assigned) return undefined
	if (future) return "upcoming"
	if (!checked || realized === undefined) return "pending"
	return realized > 0 ? "done" : "missed"
}

// `baker`'s rights at one level, undefined when it has none
export function bakerRightsAt(r: BlockRights, baker: string, currentLevel: number): BakerLevelRights | undefined {
	const [blocks = 0, attestations = 0, baked, attested] = r.rights[baker] ?? []
	if (!blocks && !attestations) return undefined
	const future = r.level > currentLevel
	return {
		level: r.level,
		blocks,
		attestations,
		blockStatus: rightStatus(blocks, future, r.realized_checked, baked),
		attestationStatus: rightStatus(attestations, future, r.realized_checked, attested),
	}
}

// levels where `baker` has rights, ascending
export function bakerRights(rights: Array<BlockRights>, baker: string, currentLevel: number) {
	return rights
		.map((r) => bakerRightsAt(r, baker, currentLevel))
		.filter((r) => r !== undefined)
		.sort((a, b) => a.level - b.level)
}

// counts over checked past levels and the nearest upcoming block
export function summarizeRights(rows: Array<BakerLevelRights>) {
	const count = (key: "blockStatus" | "attestationStatus", status: RightStatus) => rows.filter((r) => r[key] === status).length
	return {
		blocksDone: count("blockStatus", "done"),
		blocksMissed: count("blockStatus", "missed"),
		attestationsDone: count("attestationStatus", "done"),
		attestationsMissed: count("attestationStatus", "missed"),
		nextBlock: rows.find((r) => r.blockStatus === "upcoming"),
	}
}

const STATUS_PRIORITY: Array<RightStatus> = ["missed", "pending", "upcoming", "done"]

// one status for a level: a miss wins, then pending, then upcoming
export function worstStatus(r: BakerLevelRights): RightStatus {
	return STATUS_PRIORITY.find((s) => r.blockStatus === s || r.attestationStatus === s) ?? "done"
}

export function describeRights(r: BakerLevelRights) {
	const parts = []
	if (r.blocks) parts.push(`block ${r.blockStatus === "done" ? "baked" : r.blockStatus}`)
	if (r.attestations) parts.push(`${r.attestations} attestation slots ${r.attestationStatus === "done" ? "attested" : r.attestationStatus}`)
	return parts.join(", ")
}
