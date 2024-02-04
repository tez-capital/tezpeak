import { readdir, readFile, writeFile, unlink } from "fs/promises"
import { join } from "path"
import * as url from "url"

const __dirname = url.fileURLToPath(new URL('.', import.meta.url));

const sourceDir = join(__dirname, "./src")
const targetDir = join(__dirname, "./components")
for (const residue of await readdir(targetDir)) {
	unlink(join(targetDir, residue))
}
for (const file of await readdir(sourceDir)) {
	if (file.endsWith(".svg")) {
		const iconContent = await readFile(join(sourceDir, file))
		const target = file.substring(0, file.length - 4) + ".svelte"
		console.log(`writing ${target}...`)
		await writeFile(join(targetDir, target), iconContent)
	}
}
