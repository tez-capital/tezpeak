import { readable } from "svelte/store"

// ticks once a second while subscribed; used to refresh relative times ("3 seconds ago")
export const now = readable(Date.now(), (set) => {
	const interval = setInterval(() => set(Date.now()), 1000)
	return () => clearInterval(interval)
})
