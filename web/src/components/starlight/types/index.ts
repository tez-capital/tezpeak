import type { SvelteComponent } from "svelte"

export interface FontIcon {
	id: string
	family: string
	kind: "font"
}

export interface ComponentIcon {
	component: typeof SvelteComponent
	kind: "component"
}
export type IIcon = FontIcon | ComponentIcon

export type NavMenuItemType = {
	icon: IIcon
	label: string
	hidden?: boolean
	disabled?: boolean
	path?: string
	click: () => void
	isActive?: (path: string) => boolean | boolean
}

export interface MenuItem<T extends any = any> {
	icon?: IIcon
	label: string
	value: T
	hidden?: boolean
	disabled?: boolean
	isActive?: boolean
	click?: () => void
}

export type SelectItem = MenuItem

export const ComponentIcon = (component: typeof SvelteComponent): ComponentIcon => {
	return {
		component,
		kind: "component"
	}
}

export type ValidationRule = (value: any) => true | string
export type ValidationRules = Array<ValidationRule>

export type ProgressValue = number | "indeterminate"

export type PromiseFinalizers = {
	resolve: (value: any) => void
	reject: (reason?: any) => void
}

export interface DataTableItem {
	[key: string]: any,
	selected?: boolean
}

export type DataTableColumn = {
	hidden?: boolean
	name: string
	label?: string
	sortable?: boolean
	filterable?: boolean
	sortFn?: (a: DataTableItem, b: DataTableItem) => number
	getValue?: (item: DataTableItem) => any
	component?: typeof SvelteComponent<{ item: DataTableItem, compactModeEnabled: boolean }>
	click?: (item: DataTableItem) => any
	disableSelect?: boolean
}

export type DataTableFilterFunction = (item: DataTableItem, filter: any) => boolean
export enum DataTableSelectMode {
	single = "single",
	multiple = "multiple",
	none = "none"
}
export type DataTableSort = {
	column: string
	direction: "asc" | "desc"
	sortFn?: (a: DataTableItem, b: DataTableItem) => number
}
export enum DataTableHighlightMode {
	none = "none",
	row = "row",
	column = "column",
	cell = "cell",
	"row-column" = "row-column"
}
export enum DataTableLayoutMode {
	auto = "auto",
	compact = "compact",
	standard = "standard"
}

export type DataTableRowClickHandler = (item: DataTableItem) => void
export type DataTableCellClickHandler = (item: DataTableItem, column: DataTableColumn) => void

export enum ExpandingContainerDirection {
	horizontal = "horizontal",
	vertical = "vertical"
}