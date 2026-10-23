// `pnpm demo`: serves /api with generated demo data from the vite dev server, so the UI runs without a backend.
import type { IncomingMessage, ServerResponse } from 'node:http';
import type { Plugin } from 'vite';

const BAKER = 'tz1P6WKJu2rcbxKiKRZHKQKmKrpC9TfW1AwM';
const BAKER2 = 'tz1hZvgjekGo7DmQjWh7XnY5eLQD8wNYPczE';
// most operators run a single baker; the footer switch (GET /api/demo?bakers=1|2) or DEMO_BAKERS=2 adds a second one
const demo = { bakers: process.env.DEMO_BAKERS === '2' ? 2 : 1 };
const tracked = <T>(byBaker: Record<string, T>) => (demo.bakers > 1 ? byBaker : { [BAKER]: byBaker[BAKER] });
// open status streams, pushed to right away when the demo changes
const streams = new Set<() => void>();
const PAYOUT_WALLET = 'tz1gXWW1q8NcXtVy2oVVcc2s4XKNzv9CryWd';
const BLOCK_TIME = 6_000;
const CYCLE_BLOCKS = 10_800;
const START_LEVEL = 9_123_456;
const startedAt = Date.now();
const servicesStarted = new Date(startedAt - 3 * 86_400_000).toISOString();

const continual = { enabled: true, running: true };

const level = () => START_LEVEL + Math.floor((Date.now() - startedAt) / BLOCK_TIME);
const cycle = (l = level()) => Math.floor(l / CYCLE_BLOCKS);
const fakeHash = (prefix: string, seed: number, length = 51) => {
	const alphabet = '123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz';
	let out = prefix;
	for (let i = 0; out.length < length; i++) out += alphabet[(seed * 7919 + i * 104729) % alphabet.length];
	return out;
};
const address = (i: number) => fakeHash('tz1', i, 36);

function rights(current: number) {
	const result = [];
	// same window as the backend default (rights_block_window 50)
	for (let l = current - 24; l <= current + 25; l++) {
		const past = l <= current;
		// block rights are rare on mainnet (a few per cycle); keep a couple inside the window
		const blocks = l % 23 === 0 ? 1 : 0;
		const blocks2 = l % 31 === 0 ? 1 : 0;
		result.push({
			level: l,
			realized_checked: past,
			rights: tracked({
				[BAKER]: past ? [blocks, 14, blocks, l % 13 === 0 ? 0 : 1] : [blocks, 14],
				[BAKER2]: past ? [blocks2, 3, blocks2, 1] : [blocks2, 3]
			})
		});
	}
	return result;
}

function status() {
	const l = level();
	const service = { status: 'running', started: servicesStarted };
	return {
		id: 'demo',
		nodes: {
			'tezos-node': {
				address: 'http://127.0.0.1:8732',
				connection_status: 'connected',
				is_essential: true,
				network_info: { connection_count: 42 },
				block: {
					hash: fakeHash('BL', l),
					timestamp: new Date(startedAt + (l - START_LEVEL) * BLOCK_TIME).toISOString(),
					protocol: 'PsQuebecnLByd3JwTiGadoG4nGWi3HYiLXUjkibeFV8dCFeVMUg',
					level_info: { level: l, level_position: l, cycle: cycle(l), cycle_position: l % CYCLE_BLOCKS },
					voting_period_info: {
						position: 51_000,
						remaining: 24_000,
						voting_period: { index: 150, kind: 'exploration', start_position: l - 51_000 }
					}
				}
			},
			'backup-node': { address: 'http://10.0.0.2:8732', connection_status: 'disconnected', is_essential: false }
		},
		modules: {
			tezbake: {
				rights: { level: l, rights: rights(l) },
				services: {
					timestamp: Math.floor(Date.now() / 1000),
					applications: {
						node: { node: service },
						signer: { signer: service },
						baker: { baker: service, accuser: service, 'dal-node': service }
					}
				},
				bakers: {
					level: l,
					bakers: tracked({
						[BAKER]: {
							deactivated: false,
							balance: '412345678901',
							staked_balance: '400000000000',
							external_staked_balance: '125000000000',
							external_delegated_balance: '3456789000000',
							delegators_count: 1284
						},
						[BAKER2]: {
							deactivated: false,
							balance: '6000000000',
							staked_balance: '6000000000',
							external_staked_balance: '0',
							external_delegated_balance: '42000000000',
							delegators_count: 17
						}
					})
				},
				wallets: {
					baker: { app_version: '2.4.7', authorized: true, ledger: 'Nano S Plus', ledger_status: 'connected', pkh: BAKER }
				}
			},
			tezpay: {
				services: {
					timestamp: Math.floor(Date.now() / 1000),
					applications: {
						tezpay: continual.enabled
							? { 'tezpay-continual': { status: continual.running ? 'running' : 'stopped', started: servicesStarted } }
							: {}
					}
				},
				wallet: { address: PAYOUT_WALLET, balance: 842_500_000, level: 'ok' }
			}
		}
	};
}

const recipe = (i: number, c: number) => {
	const balance = 1_000_000_000 + ((i * 7_919_000_000) % 90_000_000_000);
	const valid = i % 9 !== 0;
	const amount = Math.floor(balance * 0.0009);
	return {
		baker: BAKER,
		delegator: address(i),
		cycle: c,
		recipient: address(i),
		kind: 'delegator_reward',
		fa_token_id: '0',
		fa_contract: '',
		delegator_balance: String(balance),
		amount: String(valid ? amount : 0),
		fee_rate: 0.05,
		fee: String(Math.floor(amount * 0.05)),
		note: valid ? '' : 'BALANCE_BELOW_MINIMUM',
		valid
	};
};

const summary = (c: number) => ({
	cycle: c,
	delegators: 24,
	paid_delegators: 21,
	own_staked_balance: '400000000000',
	own_delegated_balance: '12345678901',
	external_staked_balance: '125000000000',
	external_delegated_balance: '3456789000000',
	cycle_fees: '52000',
	cycle_rewards: '3215000000',
	distributed_rewards: '2700000000',
	bond_income: '380000000',
	fee_income: '135000000',
	total_income: '515000000',
	donated_bonds: '0',
	donated_fees: '6750000',
	donated_total: '6750000',
	timestamp: new Date().toISOString()
});

const report = (id: string) => {
	const c = Number(id) || cycle() - 1;
	const recipes = Array.from({ length: 24 }, (_, i) => recipe(i + 1, c));
	const toReport = (r: ReturnType<typeof recipe>, i: number) => ({
		...r,
		id: String(i),
		timestamp: new Date().toISOString(),
		tx_kind: 'tez',
		contract: '',
		token_id: '0',
		tx_fee: '1350',
		op_hash: r.valid ? fakeHash('oo', c * 100 + Math.floor(i / 10)) : '',
		success: r.valid
	});
	return {
		name: id,
		summary: summary(c),
		payouts: recipes.filter((r) => r.valid).map(toReport),
		invalid: recipes.filter((r) => !r.valid).map(toReport)
	};
};

const log = (msg: string, extra: object = {}, level = 'info') =>
	JSON.stringify({ time: new Date().toISOString(), level, msg, ...extra });

// writes newline-delimited JSON lines with a delay between them, like tezpay's streamed output
async function stream(res: ServerResponse, lines: Array<string>, delay = 400) {
	res.writeHead(200, { 'Content-Type': 'text/plain' });
	for (const line of lines) {
		res.write(line + '\n');
		await new Promise((r) => setTimeout(r, delay));
	}
	res.end(JSON.stringify({ phase: 'execution_finished', exit_code: 0, success: true }) + '\n');
}

async function readBody(req: IncomingMessage) {
	let body = '';
	for await (const chunk of req) body += chunk;
	return body;
}

export function demoApi(): Plugin {
	return {
		name: 'tezpeak-demo-api',
		configureServer(server) {
			server.middlewares.use(async (req, res, next) => {
				const url = new URL(req.url ?? '/', 'http://demo');
				if (!url.pathname.startsWith('/api/')) return next();
				const json = (value: unknown) => {
					res.writeHead(200, { 'Content-Type': 'application/json' });
					res.end(JSON.stringify(value));
				};
				const c = cycle() - 1;

				switch (url.pathname) {
					case '/api/sse': {
						res.writeHead(200, { 'Content-Type': 'text/event-stream', 'Cache-Control': 'no-cache', Connection: 'keep-alive' });
						const send = () => res.write(`data: ${JSON.stringify(status())}\n\n`);
						send();
						const timer = setInterval(send, 2_000);
						streams.add(send);
						req.on('close', () => {
							clearInterval(timer);
							streams.delete(send);
						});
						return;
					}
					case '/api/demo':
						demo.bakers = url.searchParams.get('bakers') === '2' ? 2 : 1;
						streams.forEach((send) => send());
						return json(demo);
					case '/api/tezpay/info':
						return json({
							version: { 'ami-tezpay': '0.12.0', tezpay: '0.24.1' },
							payout_wallet: PAYOUT_WALLET,
							configuration: {
								notifications: [{ type: 'discord' }, { type: 'telegram' }, { type: 'email' }],
								extensions: [{ id: 'tdc', kind: 'external' }]
							}
						});
					case '/api/tezpay/enable-continual':
						continual.enabled = continual.running = true;
						return json(true);
					case '/api/tezpay/disable-continual':
						continual.enabled = continual.running = false;
						return json(true);
					case '/api/tezpay/start-continual':
						continual.running = true;
						return json(true);
					case '/api/tezpay/stop-continual':
						continual.running = false;
						return json(true);
					case '/api/tezpay/list-reports':
						return json(url.searchParams.get('dry') === 'true' ? [String(c)] : [c, c - 1, c - 2, c - 3].map(String));
					case '/api/tezpay/report':
						return json(report(url.searchParams.get('id') ?? String(c)));
					case '/api/tezpay/generate-payouts': {
						const target = Number(url.searchParams.get('cycle')) || c;
						return stream(res, [
							log('acquiring lock', { phase: 'acquiring_lock' }),
							log('generating payout candidates', { phase: 'generate_payout_candidates', cycle: target }),
							log('collecting transaction fees', { phase: 'collect_transaction_fees' }),
							log('creating blueprint', { phase: 'create_blueprint' }),
							log('payouts generated', {
								phase: 'result',
								cycle_payout_blueprint: {
									cycles: [target],
									summary: summary(target),
									payouts: Array.from({ length: 24 }, (_, i) => recipe(i + 1, target))
								}
							})
						]);
					}
					case '/api/tezpay/pay': {
						const blueprint = JSON.parse((await readBody(req)) || '{}');
						const valid = (blueprint.payouts ?? []).filter((p: { valid: boolean }) => p.valid);
						const batches = [];
						for (let i = 0; i * 10 < valid.length; i++) batches.push(valid.slice(i * 10, i * 10 + 10));
						return stream(
							res,
							[
								log('splitting into batches', { phase: 'split_into_batches', batches: batches.length }),
								...batches.flatMap((recipes, i) => {
									const id = `batch-${i + 1}`;
									const op_hash = fakeHash('oo', Date.now() % 100_000 + i);
									return [
										log(`executing ${id}`, { phase: 'executing_batch', batch_id: id, recipes }),
										log(`${id} waiting for confirmation`, { phase: 'batch_waiting_for_confirmation', batch_id: id, op_hash }),
										log(`${id} applied`, { phase: 'batch_execution_finished', batch_id: id, op_hash })
									];
								})
							],
							700
						);
					}
					case '/api/tezpay/test-notify':
						return stream(res, [
							log('sending test notification', { notificator: url.searchParams.get('notificator') ?? 'all' }),
							log('notification sent')
						]);
					case '/api/tezpay/test-extensions':
						return stream(res, [log('loading extension tdc'), log('extension tdc responded'), log('all extensions ok')]);
					case '/api/governance/can-vote':
						return json(true);
					case '/api/governance/available-pkhs':
						return json([BAKER, BAKER2]);
					case '/api/governance/period-detail':
						return json({
							info: status().nodes['tezos-node'].block?.voting_period_info,
							voters: [
								{ pkh: BAKER, voting_power: 400_000_000_000 },
								{ pkh: BAKER2, voting_power: 6_000_000_000 },
								{ pkh: address(999), voting_power: 1_000_000_000_000 }
							],
							summary: { yay: '820000000000', nay: '60000000000', pass: '45000000000' },
							proposal: 'PsRiotumaAMotcRoDWW1bysEhQy2n1M5fy8JgRp8jjRfHGmfeA7',
							quorum: 4200,
							ballots: [{ pkh: BAKER2, ballot: 'yay' }]
						});
					case '/api/governance/upvote':
						return json(fakeHash('oo', Date.now() % 100_000));
					case '/api/governance/wait-for-apply':
						await new Promise((r) => setTimeout(r, 3_000));
						return json(true);
				}
				res.writeHead(404);
				res.end(`demo: ${url.pathname} is not mocked`);
			});
		}
	};
}
