<script lang="ts">
	import '@xterm/xterm/css/xterm.css';
	import { Terminal, type ITerminalInitOnlyOptions, type ITerminalOptions } from '@xterm/xterm';
	import { FitAddon } from '@xterm/addon-fit';

	import { onDestroy, onMount } from 'svelte';

	let terminalContainer: HTMLDivElement;
	let term: Terminal;
	let resizeObserver: ResizeObserver | undefined = undefined;
	export let options: ITerminalOptions | ITerminalInitOnlyOptions | undefined = undefined

	onMount(() => {
		term = new Terminal(options);
		term.open(terminalContainer);
		const fitAddon = new FitAddon();
		term.loadAddon(fitAddon);
		fitAddon.fit();

		resizeObserver = new ResizeObserver(() => {
			fitAddon.fit();
		});
		resizeObserver.observe(terminalContainer);
	});

	onDestroy(() => {
		term?.dispose();
		resizeObserver?.disconnect();
	});

	export function write(text: string) {
		term?.write(text);
	}

	export function clear() {
		term?.clear();
	}
</script>

<div class="terminal" bind:this={terminalContainer}></div>

<style lang="sass">
.terminal
	width: inherit
	height: inherit
</style>
