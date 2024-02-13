import type { BallotVote } from "@src/common/types/governance";
import axios, { AxiosError } from "axios";

export async function upvote_proposal(pkh: string, period: number, proposals: string[]): Promise<string | Error> {
	try {
		const response = await axios.post('/api/governance/upvote', {
			proposals: proposals,
			source: pkh,
			period: period
		});

		const opHash = response.data;
		return opHash;
	} catch (err) {
		if (err instanceof AxiosError) {
			return new Error(err.response?.data || err.message);
		}
		return err as Error;
	}
}

export async function cast_vote(pkh: string, period: number, proposal: string, ballot: BallotVote): Promise<string | Error> {
	try {
		const response = await axios.post('/api/governance/upvote', {
			proposal: proposal,
			source: pkh,
			period: period,
			ballot: ballot
		});

		const opHash = response.data as string;
		return opHash;
	} catch (err) {
		if (err instanceof AxiosError) {
			return new Error(err.response?.data || err.message);
		}
		return err as Error;
	}
}


export function openTezgov() {
	window.open('https://gov.tez.capital', '_blank');
}

export async function waitConfirmation(opHash: string) {
	try {
		const response = await axios.post('/api/governance/wait-for-apply', opHash);

		return response.data as boolean;
	} catch (err) {
		if (err instanceof AxiosError) {
			return new Error(err.response?.data || err.message);
		}
		return err as Error;
	}
}