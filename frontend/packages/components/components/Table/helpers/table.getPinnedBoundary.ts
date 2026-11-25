import { TablePin } from '../constants/table.type';
import type { TableColumnLayout } from '../types/table.row';

export const getTablePinnedBoundary = <T,>(layout: readonly TableColumnLayout<T>[], index: number) => {
  const entry = layout[index];
  if (entry?.pin === TablePin.LEFT && layout[index + 1]?.pin !== TablePin.LEFT) return TablePin.LEFT;
  if (entry?.pin === TablePin.RIGHT && layout[index - 1]?.pin !== TablePin.RIGHT) return TablePin.RIGHT;
  return TablePin.NONE;
};
