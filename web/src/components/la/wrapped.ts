import { mapValues } from "lodash-es"
import { ComponentIcon } from "@src/components/starlight/types"
import type { SvelteComponent } from "svelte"

import * as icons from "."

const wrappedIcons = mapValues(icons, (icon) => ComponentIcon(icon as (typeof SvelteComponent)))
export default wrappedIcons