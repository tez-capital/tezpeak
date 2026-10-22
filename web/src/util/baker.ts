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

// levels where `baker` has rights, split around `currentLevel`: upcoming ascending, past descending
export function bakerRightsTimeline(rights: Array<BlockRights>, baker: string, currentLevel: number) {
	const upcoming = [] as Array<BakerLevelRights>
	const past = [] as Array<BakerLevelRights>
	for (const r of rights) {
		const [blocks = 0, attestations = 0, baked, attested] = r.rights[baker] ?? []
		if (!blocks && !attestations) continue
		const future = r.level > currentLevel
		const row = {
			level: r.level,
			blocks,
			attestations,
			blockStatus: rightStatus(blocks, future, r.realized_checked, baked),
			attestationStatus: rightStatus(attestations, future, r.realized_checked, attested),
		}
		;(future ? upcoming : past).push(row)
	}
	upcoming.sort((a, b) => a.level - b.level)
	past.sort((a, b) => b.level - a.level)
	return { upcoming, past }
}

// counts over the checked past levels
export function summarizeRights(past: Array<BakerLevelRights>) {
	const count = (key: "blockStatus" | "attestationStatus", status: RightStatus) => past.filter((r) => r[key] === status).length
	return {
		blocksDone: count("blockStatus", "done"),
		blocksMissed: count("blockStatus", "missed"),
		attestationsDone: count("attestationStatus", "done"),
		attestationsMissed: count("attestationStatus", "missed"),
	}
}

const STATUS_PRIORITY: Array<RightStatus> = ["missed", "pending", "upcoming", "done"]
const worst = (a?: RightStatus, b?: RightStatus) =>
	!a ? b : !b ? a : STATUS_PRIORITY[Math.min(STATUS_PRIORITY.indexOf(a), STATUS_PRIORITY.indexOf(b))]

// rights of all `bakers` folded into one row per level (counts summed, worst status wins), same order as input
export function mergeByLevel(rows: Array<BakerLevelRights>) {
	const byLevel = new Map<number, BakerLevelRights>()
	for (const r of rows) {
		const m = byLevel.get(r.level)
		byLevel.set(r.level, !m ? { ...r } : {
			level: r.level,
			blocks: m.blocks + r.blocks,
			attestations: m.attestations + r.attestations,
			blockStatus: worst(m.blockStatus, r.blockStatus),
			attestationStatus: worst(m.attestationStatus, r.attestationStatus),
		})
	}
	return [...byLevel.values()]
}
