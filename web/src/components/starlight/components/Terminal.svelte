<script lang="ts">
	import { Terminal, type ITerminalInitOnlyOptions, type ITerminalOptions } from '@xterm/xterm';
	import { FitAddon } from '@xterm/addon-fit';

	import { onDestroy, onMount } from 'svelte';

	let terminalContainer: HTMLDivElement;
	let term: Terminal;
	let resizeObserver: ResizeObserver | undefined = undefined;
	export let options: ITerminalOptions | ITerminalInitOnlyOptions | undefined = {
		convertEol: true
	};

	onMount(() => {
		term = new Terminal(options);
		const fitAddon = new FitAddon();
		term.loadAddon(fitAddon);
		term.open(terminalContainer);
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

<div class="terminal-warp">
	<div class="terminal" bind:this={terminalContainer}></div>
</div>

<style lang="sass">
.terminal-wrap
	width: 50vw
	height: 50vh
</style>
