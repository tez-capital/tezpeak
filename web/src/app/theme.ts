import { writable } from "svelte/store"

// @alis-is/starlight themes (toned down in styles/themes.sass), applied as the only class on <html>
export const THEMES = ["dark", "light"] as const
export type Theme = (typeof THEMES)[number]

// only an explicit choice is stored; without one the theme follows the system preference
const KEY = "tezpeak-theme"
const systemLight = window.matchMedia("(prefers-color-scheme: light)")

function storedTheme(): Theme | undefined {
	try {
		const stored = localStorage.getItem(KEY) as Theme
		if (THEMES.includes(stored)) return stored
	} catch {
		// storage unavailable (private mode etc.)
	}
	return undefined
}

const systemTheme = (): Theme => (systemLight.matches ? "light" : "dark")

export const theme = writable<Theme>(storedTheme() ?? systemTheme())
theme.subscribe((t) => (document.documentElement.className = t))
systemLight.addEventListener("change", () => {
	if (!storedTheme()) theme.set(systemTheme())
})

export function cycleTheme() {
	theme.update((t) => {
		const next = THEMES[(THEMES.indexOf(t) + 1) % THEMES.length]
		try {
			localStorage.setItem(KEY, next)
		} catch {
			// storage unavailable: the choice lasts for this page only
		}
		return next
	})
}
