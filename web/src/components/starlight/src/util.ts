import type { ValidationRules } from "../types";
import { MENU_LAYER_ID } from "./constants";

export function teleportUnderElement(element: HTMLElement & { __origin?: HTMLElement | null }, targetId: string) {
	element.__origin = element.parentElement;
	element.parentElement?.removeChild(element);
	document.getElementById(targetId)?.appendChild(element);
}

export function teleportBack(element: HTMLElement & { __origin?: HTMLElement | null }) {
	element.parentElement?.removeChild(element);
	element.__origin?.appendChild(element);
	element.__origin = null;
}

export function teleportToMenulayer(element: HTMLElement & { __origin?: HTMLElement | null }) {
	teleportUnderElement(element, MENU_LAYER_ID);
}

export function createMenuLayer() {
	if (document.getElementById(MENU_LAYER_ID)) return
	const menuLayer = document.createElement("div");
	menuLayer.id = MENU_LAYER_ID;
	document.body.appendChild(menuLayer);
}

export function validate(value: any, rules: ValidationRules): true | string {
	for (const rule of rules) {
		const result = rule(value);
		if (result !== true) return result;
	}
	return true;
}

export function createKeyPressHandler(keys: Array<string>, handler: (...args: Array<any>) => any) {
	return function (event: KeyboardEvent) {
		if (keys.includes(event.key)) {
			handler(event);
		}
	}
}