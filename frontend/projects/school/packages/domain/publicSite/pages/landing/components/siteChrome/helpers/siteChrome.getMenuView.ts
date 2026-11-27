import { Menu, X } from '@Components/icons';
/** Derives the accessible mobile navigation control from open state. */
export const getSiteChromeMenuView = (open: boolean) => {
  if (open) return { menuLabel: __('Close menu'), MenuIcon: X };
  return { menuLabel: __('Open menu'), MenuIcon: Menu };
};
