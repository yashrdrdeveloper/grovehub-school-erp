import type { GetShellClassName } from './types/shell.getClassName';
export const getShellClassName: GetShellClassName = (collapsed) => {
  if (collapsed) return 'shell shell--collapsed';
  return 'shell';
};
