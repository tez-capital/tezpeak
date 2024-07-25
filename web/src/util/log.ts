import { omit } from "lodash-es"
import { stringify } from "hjson"

export type LogMessage = object & { time: string, level?: string, msg: string, error?: string, phase?: string }

export function formatLogMessageForTerminal(message: string | LogMessage): string {
	if (typeof message === "string") {
		message = JSON.parse(message) as LogMessage
	}

	if (message.phase === "execution_finished") {
		return "==== execution finished ====\n"
	}

	let levelString = `[${message.level}]`
	switch (message.level?.toLowerCase()) {
		case "error":
			levelString = `\x1b[31m[${message.level}]\x1b[0m`
			break
		case "warn":
		case "warning":
			levelString = `\x1b[33m[${message.level}]\x1b[0m`
			break
	}
	
	const time = new Date(message.time).toLocaleString()

	let line = `${time} ${levelString} ${message.msg}\n`
	const restOfProperties = omit(message, ["time", "level", "msg", "error"])
	if (Object.keys(restOfProperties).length > 0) {
		line += `${stringify(restOfProperties, { separator: true, space: "\t" })}\n`
	}
	if (message.error) {
		line += `\x1b[31m[ERROR]\x1b[0m: ${message.error.replace("\\n", "\n")}\n` // Replace \n with actual new line
	}
	return line
}