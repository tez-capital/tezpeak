export async function readBody(response: Response, cb: (message: string) => void) {
	const reader = response.body?.getReader();
	if (!reader) {
		throw new Error('reader not available');
	}

	const decoder = new TextDecoder('utf-8');
	let buffer = '';

	for (; ;) {
		const { done, value } = await reader.read();
		if (done) {
			if (buffer.length > 0) {
				console.log(`Received chunk: ${buffer}`);
			}
			break;
		}

		buffer += decoder.decode(value, { stream: true });

		let newlineIndex;
		while ((newlineIndex = buffer.indexOf('\n')) > -1) {
			const line = buffer.slice(0, newlineIndex + 1).trim();
			buffer = buffer.slice(newlineIndex + 1);
			cb(line);
		}
	}
}