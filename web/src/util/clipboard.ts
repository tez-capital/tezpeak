export async function writeToClipboard(data: string) {
	await navigator.clipboard.writeText(data)
}