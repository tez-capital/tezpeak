<script lang="ts">
	import Select from '@components/starlight/components/Select.svelte';
	import {
		DataTableSelectMode,
		type DataTableColumn,
		type DataTableFilterFunction,
		type DataTableItem,
		type DataTableSort,
		DataTableLayoutMode,
		DataTableHighlightMode,
		type DataTableRowClickHandler,
		type DataTableCellClickHandler
	} from '../types';
	import CheckBox from './CheckBox.svelte';
	import ArrowDown from '../icons/arrow-down.svelte';
	import { onDestroy, onMount } from 'svelte';
	import { createKeyPressHandler } from '../src/util';

	export let hideControls = false;
	export let compactMode = false;
	export let selectMode: DataTableSelectMode = DataTableSelectMode.none;
	$: selectionEnabled = selectMode !== DataTableSelectMode.none;
	export let layoutMode: DataTableLayoutMode = DataTableLayoutMode.standard;
	export let highlightMode: DataTableHighlightMode = DataTableHighlightMode.none;
	export let enableRowClickSelect = false;

	let table: HTMLDivElement;

	export let columns: Array<DataTableColumn> = [];
	$: dataColumns =
		columns && columns.length > 0
			? columns.filter((column) => !column.hidden)
			: Object.keys(preprocessedData[0] ?? {}).map(
					(key) => ({ name: key, label: key, sortable: true, filterable: true }) as DataTableColumn
			  );

	//#region data preprocessing
	export let filter = '';
	export let filterFunction: DataTableFilterFunction = (item, filter: string) => {
		if (filter === '') {
			return true;
		}
		for (const column of dataColumns) {
			if (column.filterable) {
				let value = item[column.name];
				if (column.getValue !== undefined) {
					value = column.getValue(item);
				}
				if (value !== undefined && value !== null) {
					if (value.toString().toLowerCase().includes(filter.toLowerCase())) {
						return true;
					}
				}
			}
		}
		return false;
	};

	export let activeSort: DataTableSort | undefined = undefined;
	function defaultSortFn(a: DataTableItem, b: DataTableItem) {
		if (!activeSort) {
			return 0;
		}
		const aValue = a[activeSort.column];
		const bValue = b[activeSort.column];
		if (aValue === undefined || aValue === null) return 1;
		if (bValue === undefined || bValue === null) return -1;

		if (activeSort.direction === 'asc') return aValue > bValue ? 1 : -1;

		return aValue < bValue ? 1 : -1;
	}

	export let data: Array<DataTableItem>;
	$: preprocessedData = data.filter((item) => filterFunction(item, filter));
	$: dataToDisplay = !activeSort
		? preprocessedData
		: preprocessedData.sort(activeSort?.sortFn ?? defaultSortFn);
	//#endregion

	//#region selection
	function selectItem(item: DataTableItem) {
		if (selectMode === DataTableSelectMode.single) {
			for (const otherItem of data) {
				otherItem.selected = false;
			}
		}
		item.selected = !item.selected;
		data = [...data];
	}
	function selectAll() {
		for (const item of data) {
			item.selected = !allSelected;
		}
		data = [...data];
	}

	$: allSelected = data.every((item) => item.selected);
	//#endregion

	//#region automatic resizing
	let widthObserverTable: HTMLTableElement;
	let widthObserverHeader: HTMLTableRowElement;
	let compactModeColumnCount = 1;
	let compactModeRequired = false;
	export let optimalCompactModePropertyWidth = 200;
	function updateCompactModeColumnCount() {
		compactModeColumnCount = Math.max(
			Math.floor(table.clientWidth / optimalCompactModePropertyWidth),
			1
		);
	}

	function updateLayout() {
		switch (layoutMode) {
			case DataTableLayoutMode.auto:
				if (!widthObserverTable || !widthObserverTable.parentElement || !widthObserverHeader) {
					return;
				}
				if (widthObserverTable.parentElement.clientWidth > widthObserverTable.clientWidth) return;

				const spaceFills = widthObserverHeader.querySelectorAll('.space-fill');
				for (const spaceFill of spaceFills) {
					if (spaceFill.clientWidth === 0) {
						compactModeRequired = true;
						updateCompactModeColumnCount();
						return;
					}
				}
				compactModeRequired = false;
				break;
			case DataTableLayoutMode.compact:
				updateCompactModeColumnCount();
				break;
			default:
				return;
		}
	}

	let widthDebounceTimer: NodeJS.Timeout | undefined;
	const widthResizeObserver = new ResizeObserver(async () => {
		clearTimeout(widthDebounceTimer);
		widthDebounceTimer = setTimeout(updateLayout, 10);
	});

	$: compactModeEnabled = layoutMode === DataTableLayoutMode.compact || compactModeRequired;
	//#endregion

	//#region higlighting
	let highlightMaxWidth = 0;
	let highlightMaxHeight = 0;
	function updateHighlightLimits() {
		highlightMaxWidth = table?.scrollWidth;
		highlightMaxHeight = table?.scrollHeight;
	}
	let debounceSourceUpdate: NodeJS.Timeout | undefined;
	type HiglightIgnore = 'column' | 'row' | undefined;
	type HighlightInfo = {
		left: number;
		top: number;
		width: number;
		height: number;
		highlight: boolean;
		ignore?: HiglightIgnore;
	};
	let highlightInfo: HighlightInfo = { left: 0, top: 0, width: 0, height: 0, highlight: false };
	function updateHighlightSource(
		e: Event,
		options?: { ignore?: HiglightIgnore; remove?: boolean }
	) {
		clearTimeout(debounceSourceUpdate);
		debounceSourceUpdate = setTimeout(() => {
			if (options?.remove) {
				highlightInfo = { ...highlightInfo, highlight: false, ignore: undefined };
				return;
			}
			const el = e.target as HTMLElement;
			highlightInfo = {
				left: el.offsetLeft,
				top: el.offsetTop,
				width: el.clientWidth,
				height: el.clientHeight,
				highlight: true,
				ignore: options?.ignore
			};
		}, 20);
	}
	function createHighlightEventHandler(options?: { ignore?: HiglightIgnore; remove?: boolean }) {
		return (e: Event) => updateHighlightSource(e, options);
	}

	let debounceTimer: NodeJS.Timeout | undefined;
	const resizeObserver = new ResizeObserver(async () => {
		clearTimeout(debounceTimer);
		debounceTimer = setTimeout(updateHighlightLimits, 10);
	});
	//#endregion

	// common

	onMount(() => {
		updateLayout();
		updateHighlightLimits();
		widthResizeObserver.observe(widthObserverHeader);
		resizeObserver.observe(table);
	});
	onDestroy(() => {
		widthResizeObserver.disconnect();
		resizeObserver.disconnect();
	});

	function sortBy(column: DataTableColumn) {
		if (activeSort?.column === column.name) {
			activeSort.direction = activeSort.direction === 'asc' ? 'desc' : 'asc';
		} else {
			activeSort = { column: column.name, direction: 'asc', sortFn: column.sortFn };
		}
	}

	export let rowClickHandler: DataTableRowClickHandler | undefined = undefined;
	export let cellClickHandler: DataTableCellClickHandler | undefined = undefined;
</script>

<div class="data-table">
	<div class="width-observer">
		<table bind:this={widthObserverTable}>
			<thead>
			<tr bind:this={widthObserverHeader}>
				{#if selectionEnabled}
					<div class="header cell checkbox*cell">
						<CheckBox class="checkbox" />
					</div>
				{/if}
				<!-- slot header -->
				<slot name="header" {dataColumns} {activeSort} {compactModeEnabled}>
					{#each dataColumns as column}
						<th class:sortable-column={column.sortable}>
							<slot name="column-header" {column} {activeSort} {compactModeEnabled}>
								<div class="column-header-content">
									<b class="column-header-label">{column.label ?? column.name}</b>
									<div class="space-fill"></div>
									<div class="sort-arrow">
										<ArrowDown />
									</div>
								</div>
							</slot>
						</th>
					{/each}
				</slot>
			</tr>
			</thead>
		</table>
	</div>
	{#if !hideControls}
		<div class="controls">
			{#if compactMode}
				<div class="sort-select">
					<Select label="Order by" />
				</div>
			{/if}
		</div>
	{/if}
	<div
		class="table-grid"
		bind:this={table}
		style:--column-count={dataColumns.length}
		style:--row-count={Math.max(dataToDisplay.length, 1)}
		class:selectable={selectionEnabled}
		style:--highlight-max-width={`${highlightMaxWidth}px`}
		style:--highlight-max-height={`${highlightMaxHeight}px`}
		style:--highlight-source-left={`${highlightInfo.left}px`}
		style:--highlight-source-top={`${highlightInfo.top}px`}
		style:--highlight-source-width={`${highlightInfo.width}px`}
		style:--highlight-source-height={`${highlightInfo.height}px`}
		class:highlight-column={[
			DataTableHighlightMode.column,
			DataTableHighlightMode['row-column']
		].includes(highlightMode)}
		class:highlight-row={[
			DataTableHighlightMode.row,
			DataTableHighlightMode['row-column']
		].includes(highlightMode)}
		class:highlight={highlightMode !== DataTableHighlightMode.none}
		class:highlight-column-visible={highlightInfo.highlight && highlightInfo.ignore !== 'column'}
		class:highlight-row-visible={highlightInfo.highlight && highlightInfo.ignore !== 'row'}
	>
		<div class="table-header">
			{#if !compactModeEnabled}
				{@const showHighlight = createHighlightEventHandler({ ignore: 'row' })}
				{@const hideHighlight = createHighlightEventHandler({ remove: true })}
				{#if selectionEnabled}
					<button
						class="unstyle-button header cell checkbox-cell"
						on:click={() => selectMode === DataTableSelectMode.multiple && selectAll()}
						on:keypress={createKeyPressHandler(
							['Enter'],
							() => selectMode === DataTableSelectMode.multiple && selectAll()
						)}
					>
						{#if selectMode === DataTableSelectMode.multiple}
							<CheckBox class="checkbox" checked={allSelected} />
						{/if}
					</button>
				{/if}
				<slot name="header" {dataColumns} {activeSort} {compactModeEnabled}>
					{#each dataColumns as column}
						<button
							on:click={() => sortBy(column)}
							on:pointerenter={showHighlight}
							on:mouseenter={showHighlight}
							on:pointerleave={hideHighlight}
							on:mouseleave={hideHighlight}
							on:keypress={createKeyPressHandler(['Enter'], () => sortBy(column))}
							aria-label={`Sort by ${column.label ?? column.name}`}
							class="unstyle-button header cell"
							class:sortable-column={column.sortable}
						>
							<slot name="column-header" {column} {activeSort} {compactModeEnabled}>
								<div class="column-header-content">
									<b class="column-header-label">{column.label ?? column.name}</b>
									<!-- <div class="space-fill"></div> -->
									{#if column.sortable}
									<div
										class="sort-arrow"
										class:sort-arrow-active={activeSort && activeSort?.column === column.name}
										class:sort-arrow-asc={activeSort &&
											activeSort?.column === column.name &&
											activeSort?.direction === 'asc'}
									>
										<ArrowDown />
									</div>
									{/if}
								</div>
							</slot>
						</button>
					{/each}
				</slot>
			{/if}
		</div>
		<div class="table-body">
			{#if dataToDisplay && dataToDisplay.length > 0}
				{@const showHighlight = createHighlightEventHandler()}
				{@const hideHighlight = createHighlightEventHandler({ remove: true })}
				{#each dataToDisplay as item, i}
					{#if selectionEnabled}
						{@const showHighlight = createHighlightEventHandler({ ignore: 'column' })}
						<button
							class="unstyle-button cell checkbox-cell"
							class:even-row={i % 2 === 0}
							on:pointerenter={showHighlight}
							on:mouseenter={showHighlight}
							on:pointerleave={hideHighlight}
							on:mouseleave={hideHighlight}
							on:click={() => selectionEnabled && selectItem(item)}
						>
							<CheckBox class="checkbox" checked={item.selected} />
						</button>
					{/if}
					{#if !compactModeEnabled}
						{#each dataColumns as column}
							{@const columnSelectionEnabled = enableRowClickSelect && selectionEnabled && !column.disableSelect && !column.click}
							<button
								class="unstyle-button cell"
								class:even-row={i % 2 === 0}
								on:click={() => cellClickHandler && cellClickHandler(item, column)}
								on:click={() => rowClickHandler && rowClickHandler(item)}
								class:clickable={rowClickHandler !== undefined ||
									columnSelectionEnabled ||
									column.click}
								on:click={() => column.click && column.click(item)}
								on:click={() => columnSelectionEnabled && selectItem(item)}
								on:pointerenter={showHighlight}
								on:mouseenter={showHighlight}
								on:pointerleave={hideHighlight}
								on:mouseleave={hideHighlight}
								class:highlight-cell={[DataTableHighlightMode.cell].includes(highlightMode)}
							>
								{#if column.component !== undefined}
									<svelte:component this={column.component} {item} {compactModeEnabled} />
								{:else if column.getValue !== undefined}
									{column.getValue(item)}
								{:else}
									{item[column.name] ?? `-`}
								{/if}
							</button>
						{/each}
					{:else}
						{#each dataColumns as column}
							<button
								class="unstyle-button compact-mode-row"
								class:clickable={rowClickHandler !== undefined || selectionEnabled}
								on:click={() => rowClickHandler && rowClickHandler(item)}
								on:click={() => selectionEnabled && selectItem(item)}
							>
								{#if column.component !== undefined}
									<svelte:component this={column.component} {item} {compactModeEnabled} />
								{:else}
									<b class="compact-mode-row-header">{column.label ?? column.name}</b>

									<div class="compact-mode-row-value">
										{#if column.getValue !== undefined}
											{column.getValue(item)}
										{:else}
											{item[column.name] ?? `-`}
										{/if}
									</div>
								{/if}
							</button>
						{/each}
					{/if}
				{/each}
			{:else}
				<slot name="no-data">
					<div class="no-data">No data</div>
				</slot>
			{/if}
		</div>
		<div class="row-highlighter" />
		<div class="column-highlighter" />
	</div>
</div>

<!-- {JSON.stringify(sortedAndFilteredData)} -->

<style lang="sass">
.data-table
	position: relative
	width: 100%
	max-width: 100%
	height: 100%
	overflow: hidden
	--checkbox-background-color: var(--data-table-checkbox-background-color)

	.width-observer
		max-height: 0
		overflow: hidden

	.table-grid
		display: grid
		grid-template-columns: repeat(var(--column-count), 1fr)
		grid-template-rows: repeat(var(--row-count), auto) 1fr
		height: 100%
		border-collapse: collapse
		color: var(--text-color)
		border-radius: var(--border-radius)
		table-layout: var(--data-table-table-layout, auto)
		overflow: auto
		position: relative
		z-index: 2

		&.selectable
			grid-template-columns: auto repeat(var(--column-count), 1fr)

		.table-header
			display: contents

		.header
			position: sticky
			top: 0
			z-index: 10
			border-bottom: 1px solid var(--text-color)
			background-color: var(--card-background-color)
			user-select: none !important

		.table-body
			display: contents

		.checkbox-cell
			padding: 0px !important
			cursor: pointer !important
			:global(.checkbox)
				padding: var(--spacing)
				height: calc(100% - 2* var(--spacing))
				pointer-events: none

		.cell
			display: flex
			align-items: center
			white-space: nowrap
			user-select: text
			cursor: auto
			padding: var(--spacing)

		.even-row
			background-color: var(--data-table-even-row-background-color, transparent)

		.no-data
			grid-column: 1 / -1
			padding: var(--spacing)
			text-align: center

		.column-header-content
			display: grid
			width: 100%
			align-items: center
			justify-content: center
			text-transform: capitalize
			grid-template-columns: auto 1fr auto

			.column-header-label
				grid-column: 1

			.sort-arrow
				grid-column: 3


		.sortable-column
			cursor: pointer

			.sort-arrow
				display: flex
				opacity: 0
				transition-duration: 0.5s, 0.1s
				transition-property: transform, opacity
				transform: rotate(0deg)
				align-items: center

			&:hover
				.sort-arrow
					opacity: 0.7
			
			.sort-arrow-asc
				transform: rotate(180deg)

			.sort-arrow-active
				opacity: 1

		.clickable
			cursor: pointer
			transition-duration: 0.1s

		.highlight-cell
			&:hover
				background-color: var(--data-table-cell-highlight-color)

	.highlight 	
		.row-highlighter
			position: absolute
			opacity: 0
			z-index: 10
			background-color: var(--data-table-column-highlight-color)
			pointer-events: none
			border-radius: var(--border-radius)

		.column-highlighter
			position: absolute
			opacity: 0
			z-index: 10
			background-color: var(--data-table-column-highlight-color)
			pointer-events: none
			border-radius: var(--border-radius)

	.highlight-row
		.row-highlighter 
			left: 0
			top: var(--highlight-source-top)
			width: var(--highlight-max-width)
			height: var(--highlight-source-height)
		&.highlight-row-visible
			.row-highlighter
				opacity: 1

	.highlight-column
		.column-highlighter
			left: var(--highlight-source-left)
			top: 0
			width: var(--highlight-source-width)
			height: var(--highlight-max-height)
		
		&.highlight-column-visible
			.column-highlighter
				opacity: 1

.sort-arrow
	:global(svg)
		width: calc(var(--font-size) * 1.2)




</style>
