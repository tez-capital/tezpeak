import { formatDistanceStrict, formatDistance } from "date-fns"

export function formatBlockHash(hash: string) {
  if (!hash) {
    return "N/A"
  }
  return `${hash.slice(0, 10)}...${hash.slice(-10)}`
}

export function formatAddress(address: string) {
  if (!address) {
    return "N/A"
  }
  return `${address.slice(0, 10)}...${address.slice(-10)}`
}

export function formatTimestamp(timestamp: string | number) {
  if (!timestamp) {
    return "N/A"
  }
  return new Date(timestamp).toLocaleString()
}

export function formatTimestampAgo(timestamp: string | number) {
  if (!timestamp) {
    return "N/A"
  }
  return formatDistance(new Date(timestamp), new Date())
}

export function formatTimestampAgoStrict(timestamp: string | number) {
  if (!timestamp) {
    return "N/A"
  }
  return formatDistanceStrict(new Date(timestamp), new Date(), { addSuffix: true })
}

export function formatBalance(balance: string | bigint) {
  if (typeof balance === "bigint") {
    balance = balance.toString()
  }
  if (!balance || typeof balance !== "string") {
    return "N/A"
  }
  if (balance === "0") {
    return "0 ꜩ"
  }
  return `${balance.substring(0, balance.length - 6)}.${balance.substring(balance.length - 6)} ꜩ`
}