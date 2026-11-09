// Author: HyperparameterSI

/** The attribute of `<html>` that marks a page whose window is not drawn as
 * active, for CSS: `:root[data-golet-window-inactive] .sidebar { … }`. */
export const INACTIVE_ATTRIBUTE = "data-golet-window-inactive";

/** Marks root, a page's `<html>`, as of a window drawn as active or not. */
export function markWindowActive(
  root: Pick<Element, "toggleAttribute"> | null | undefined,
  active: boolean,
): void {
  root?.toggleAttribute(INACTIVE_ATTRIBUTE, !active);
}
