import { formatDistanceStrict, formatDistance } from "date-fns"
import { BLOCK_TIME } from "@src/common/constants"

export function formatBlockHash(hash: string) {
  if (!hash) {
    return "N/A"
  }
  return `${hash.slice(0, 8)}...${hash.slice(-8)}`
}

// tz1hZ…PczE, for labels next to a baker colour
export function formatShortAddress(address: string) {
  return address ? `${address.slice(0, 5)}…${address.slice(-4)}` : "N/A"
}

export function formatAddress(address: string) {
  if (!address) {
    return "N/A"
  }
  return `${address.slice(0, 8)}...${address.slice(-8)}`
}

function timestampToJSTimestamp(timestamp: number) {
  if (timestamp < 10000000000) {
    return timestamp * 1000
  }
  return timestamp
}

export function formatTimestamp(timestamp: string | number) {
  if (!timestamp) {
    return "N/A"
  }

  if (typeof timestamp === "number") timestamp = timestampToJSTimestamp(timestamp)

  return new Date(timestamp).toLocaleString()
}

export function formatTimestampAgo(timestamp: string | number) {
  if (!timestamp) {
    return "N/A"
  }

  if (typeof timestamp === "number") timestamp = timestampToJSTimestamp(timestamp)

  return formatDistance(new Date(timestamp), new Date())
}

export function formatTimestampAgoStrict(timestamp: string | number) {
  if (!timestamp) {
    return "N/A"
  }

  if (typeof timestamp === "number") timestamp = timestampToJSTimestamp(timestamp)

  return formatDistanceStrict(new Date(timestamp), new Date(), { addSuffix: true })
}

export function formatBalance(mutez: string | bigint | number) {
  if (typeof mutez === "bigint" || typeof mutez === "number") {
    mutez = mutez.toString()
  }
  if (!mutez || typeof mutez !== "string") {
    return "N/A"
  }
  if (mutez === "0") {
    return "0 ꜩ"
  }
  mutez = mutez.padStart(7, "0")
  return `${mutez.substring(0, mutez.length - 6)}.${mutez.substring(mutez.length - 6)} ꜩ`
}

export function formatPercentage(percentage: number | string) {
  return `${Number(percentage).toFixed(2)}%`
}
// tez for reading at a glance: grouped with at most 2 decimals, or compact (3.5M ꜩ); formatBalance keeps exact mutez
export function formatTez(mutez: string | number | bigint, compact = false) {
  const tez = Number(mutez) / 1e6
  if (mutez === "" || !Number.isFinite(tez)) {
    return "N/A"
  }
  const format = compact ? { notation: "compact", maximumFractionDigits: 1 } as const : { maximumFractionDigits: 2 }
  return `${new Intl.NumberFormat(undefined, format).format(tez)} ꜩ`
}

// rough estimate for `blocks` levels from now: { in: "~3 min", at: "14:52" }
export function formatBlocksEta(blocks: number, from = Date.now()) {
  const ms = blocks * BLOCK_TIME
  const minutes = Math.round(ms / 60_000)
  const duration = minutes < 1 ? "<1 min" : minutes < 60 ? `${minutes} min` : `${Math.floor(minutes / 60)} h ${minutes % 60} min`
  return { in: `~${duration}`, at: new Date(from + ms).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" }) }
}
