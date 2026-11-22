import type { GetShellNavClassName } from './types/shell.getNavClassName';
export const getShellNavClassName: GetShellNavClassName = ({ isActive }) => {
  if (isActive) return 'nav-item nav-item--active';
  return 'nav-item';
};
