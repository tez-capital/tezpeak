// @alis-is/starlight exposes these through its untyped "./*" export (TS looks for `*.d.svelte.ts`,
// svelte-package emits `*.svelte.d.ts`), so they need ambient types until upstream adds a types condition.

declare module '@alis-is/starlight/la/icons/*.svelte' {
	import type { Component } from 'svelte';
	const Icon: Component;
	export default Icon;
}

declare module '@alis-is/starlight/dialogs/Terminal.svelte' {
	import type { Component } from 'svelte';
	type TerminalDialogState = { title: string; allow_close: boolean; close_text: string };
	type TerminalDialogExports = {
		show(options: Partial<TerminalDialogState>): Promise<void>;
		update_context(options: Partial<TerminalDialogState>): Promise<void>;
		write(data: string): Promise<void>;
		clear(): Promise<void>;
		hide(): Promise<void>;
	};
	const TerminalDialog: Component<Record<string, never>, TerminalDialogExports>;
	type TerminalDialog = TerminalDialogExports;
	export default TerminalDialog;
}
