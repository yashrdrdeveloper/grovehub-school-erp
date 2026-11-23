import { ChevronsLeft, ChevronsRight } from '@Components/icons';

/** Pure sidebar-control projection, with translated literal labels. */
export const getSidebarCollapseView = (collapsed: boolean) => {
  if (collapsed) {
    return {
      collapsedLabel: __('Expand sidebar'),
      collapsedShortLabel: __('Expand'),
      CollapseIcon: ChevronsRight,
    };
  }
  return {
    collapsedLabel: __('Collapse sidebar'),
    collapsedShortLabel: __('Collapse sidebar'),
    CollapseIcon: ChevronsLeft,
  };
};
