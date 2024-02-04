import WrappedLaIcons from "@src/components/la/wrapped";
import type { NavMenuItemType } from "@src/components/starlight/types";

export const get_items = () => {
	return [
		{ label: 'Home', icon: WrappedLaIcons.HomeIcon, path: "/" },
		{ label: 'Market', icon: WrappedLaIcons.BalanceScaleIcon, path: "/market", disabled: true },
	].filter(x => !!x) as NavMenuItemType[]
}
