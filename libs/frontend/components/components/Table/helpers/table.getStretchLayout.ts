import { TablePin } from '../constants/table.type';
import type { TableColumnLayout } from '../types/table.row';
import { getTableFillWidth } from './table.getFillWidth';

export const getTableStretchLayout = <T>(
  layout: readonly TableColumnLayout<T>[],
  available: number,
  fixedWidth: number,
  widths: Readonly<Record<string, number>>,
): readonly TableColumnLayout<T>[] => {
  const total = layout.reduce((sum, entry) => sum + entry.width, 0);
  const remaining = getTableFillWidth(available, total, fixedWidth);
  if (!remaining) return layout;
  const stretch = layout.find(
    (entry) => entry.column.fillRemainingWidth && widths[entry.column.id] === undefined,
  );
  if (!stretch) return layout;
  const growth = Math.max(
    0,
    Math.min(remaining, (stretch.column.maxWidth ?? 1200) - stretch.width),
  );
  if (!growth) return layout;
  const expanded = layout.map((entry) => {
    if (entry === stretch) return { ...entry, width: entry.width + growth };
    return entry;
  });
  return expanded.map((entry, index) => {
    if (entry.pin === TablePin.LEFT)
      return {
        ...entry,
        offset: expanded
          .slice(0, index)
          .filter((item) => item.pin === TablePin.LEFT)
          .reduce((sum, item) => sum + item.width, 0),
      };
    if (entry.pin === TablePin.RIGHT)
      return {
        ...entry,
        offset: expanded
          .slice(index + 1)
          .filter((item) => item.pin === TablePin.RIGHT)
          .reduce((sum, item) => sum + item.width, 0),
      };
    return entry;
  });
};
