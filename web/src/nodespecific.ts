// import { Buffer } from "buffer";
// globalThis.Buffer = Buffer;

// crypto.randomUUID exists only in secure contexts (https/localhost); tezpeak is often served over plain http,
// and @alis-is/starlight inputs use it for element ids. ponytail: drop once starlight stops calling it.
if (typeof crypto !== "undefined" && typeof crypto.randomUUID !== "function") {
	crypto.randomUUID = () =>
		"10000000-1000-4000-8000-100000000000".replace(/[018]/g, (c) =>
			(+c ^ (crypto.getRandomValues(new Uint8Array(1))[0] & (15 >> (+c / 4)))).toString(16),
		) as `${string}-${string}-${string}-${string}-${string}`
}
