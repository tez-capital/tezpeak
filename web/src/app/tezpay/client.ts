import { EmptyTezpayInfo, type PayoutBlueprint, type TezpayInfo } from "@src/common/types/tezpay"

export async function getTezpayInfo() {
	try {
		const response = await fetch('/api/tezpay/info', {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json'
			}
		})

		if (response.status !== 200) {
			throw new Error('Failed to get tezpay info')
		}

		return await response.json() as TezpayInfo
	} catch (e) {
		console.log(e)
		return EmptyTezpayInfo
	}
}

export async function generatePayuts(cycle: number | undefined, cb: (message: string) => void, dry?: boolean) {
	const cycleQuery = cycle ? `cycle=${cycle}` : ''
	const dryQuery = dry ? `dry=${dry}` : ''

	const response = await fetch(`/api/tezpay/generate-payouts?${cycleQuery}&${dryQuery}`, {
		method: 'GET',
		headers: {
			'Content-Type': 'application/json'
		}
	})

	if (response.status !== 200) {
		throw new Error(response.statusText)
	}

	const reader = response.body?.getReader();
	if (!reader) {
		throw new Error('No reader')
	}
	const decoder = new TextDecoder('utf-8');
	let buffer = '';


	for (;;) {
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

export async function executePayuts(blueprint: PayoutBlueprint, cb: (message: string) => void, dry?: boolean) {
	const dryQuery = dry ? `dry=${dry}` : ''

	const response = await fetch(`/api/tezpay/pay?${dryQuery}`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify(blueprint)
	})
	if (response.status !== 200) {
		throw new Error(response.statusText)
	}

	const reader = response.body?.getReader();
	if (!reader) {
		throw new Error('No reader')
	}
	const decoder = new TextDecoder('utf-8');
	let buffer = '';


	for (;;) {
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

export async function stopContinual() {
	const response = await fetch('/api/tezpay/stop-continual', {
		method: 'GET',
	})

	if (response.status !== 200) {
		const body = await response.text()
		throw new Error(body)
	}
}

export async function startContinual() {
	const response = await fetch('/api/tezpay/start-continual', {
		method: 'GET',
	})

	
	if (response.status !== 200) {
		const body = await response.text()
		console.log(body)
		throw new Error(body)
	}
}

export async function listReports(dry?: boolean) {
	const reponse = await fetch(`/api/tezpay/list-reports?dry=${dry === true}`, {
		method: 'GET',
	})

	if (reponse.status !== 200) {
		throw new Error('Failed to list reports')
	}

	return await reponse.json() as Array<string>
}

export async function getReport(report: string, dry?: boolean) {
	const reponse = await fetch(`/api/tezpay/report?id=${report}&dry=${dry === true}`, {
		method: 'GET',
	})

	if (reponse.status !== 200) {
		throw new Error('Failed to get report')
	}

	return await reponse.json()
}