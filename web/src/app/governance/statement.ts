import bs58check from "bs58check"
import bs58 from "bs58"

export const ALPHABET = '123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz';
const protoDataSize = 2 + 32 // Pt + 32 bytes sha256 hash
const base58Size = 51

const b58dec = (word: string) => {
	let x = 0n;
	for (const c of word) {
		x *= 58n;
		x += BigInt(ALPHABET.indexOf(c));
	}
	return x;
}

const asciidec = (val: bigint) => {
	const word = [];
	while (val > 0n) {
		word.push(Number(val % 256n));
		val = val / 256n;
	}
	return word.reverse();
}

export const generateStatementProtocolBytes = (dataToEncode: string) => {
	const length = protoDataSize
	const target = b58dec(dataToEncode);

	const pl = bs58.decode(dataToEncode).length + 1
	const shift = BigInt(8 * (length - pl + 4));
	const m = base58Size - dataToEncode.length

	// The first bytes in base58 must be the given target. lo is the version
	// byte we want to compute.
	let lo = target * 58n ** BigInt(m);
	lo = lo >> shift;
	// (a >> b) << b == a iff a = 2**b. Then this line is equivalent to add 1
	// if lo is not a power of 2.
	lo += lo === ((lo >> shift) << shift) ? 0n : 1n;
	// hi will be negative because of the minus if m is not big enough to
	// encode a payload of [shift] bits.
	let hi = (target + 1n) * 58n ** BigInt(m) - (1n << shift) + 1n;
	// if b is negative (i.e. m not big enough to encode a payload of the
	// given length), b >> a is always strictly negative
	hi = hi >> shift;

	if (hi >= lo) {
		const to_encode = Uint8Array.from(Array.from({ ...asciidec(lo), length }).map(x => x === undefined ? 0 : x));

		const base58_encoded_minimal_value = bs58check.encode(to_encode);
		if (!base58_encoded_minimal_value.startsWith(dataToEncode) || base58_encoded_minimal_value.length !== m + dataToEncode.length) {
			return null
		}

		const to_encode_max = Uint8Array.from(Array.from({ ...asciidec(lo), length }).map(x => x === undefined ? 255 : x));
		const base58_encoded_maximal_value = bs58check.encode(to_encode_max);
		if (!base58_encoded_maximal_value.startsWith(dataToEncode) || base58_encoded_maximal_value.length !== m + dataToEncode.length) {
			return null
		}

		return to_encode;
	}

	return null;
}

export const generateStatementProtocol = (message: string): [boolean, string] => {
	let protoMessage = "Pt1"
	for (const c of message) {
		if (ALPHABET.includes(c)) {
			protoMessage += c
			continue
		}
		if (ALPHABET.includes(c.toLowerCase())) {
			protoMessage += c.toLowerCase()
			continue
		}
	}

	const proto = generateStatementProtocolBytes(protoMessage)
	if (proto === null) {
		return [false, "Can not generate protocol hash from provided message"]
	}
	return [true, bs58check.encode(new Uint8Array(proto))]
}