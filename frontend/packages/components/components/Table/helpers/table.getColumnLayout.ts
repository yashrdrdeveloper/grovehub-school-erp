import { TablePin } from '../constants/table.type';
import type { TableColumn } from '../types/table.column';
import type { TablePreferences } from '../types/table.props';
import type { TableColumnLayout } from '../types/table.row';
import { getTableColumnWidth } from './table.preferences';

export const getTableColumnLayout = <T>(
  columns: readonly TableColumn<T>[],
  preferences: TablePreferences,
): TableColumnLayout<T>[] => {
  const byId = new Map(columns.map((column) => [column.id, column]));
  const order = Array.from(new Set([...preferences.order, ...columns.map((column) => column.id)]));
  const visible: TableColumn<T>[] = [];
  for (const id of order) {
    const column = byId.get(id);
    if (column && (column.hideable === false || !preferences.hidden.includes(id)))
      visible.push(column);
  }
  if (!visible.length && columns[0]) visible.push(columns[0]);
  const entries = visible.map((column) => ({
    column,
    width: getTableColumnWidth(column, preferences.widths[column.id]),
    pin: preferences.pins[column.id] ?? column.pin ?? TablePin.NONE,
  }));
  const left = entries.filter((item) => item.pin === TablePin.LEFT);
  const middle = entries.filter((item) => item.pin === TablePin.NONE);
  const right = entries.filter((item) => item.pin === TablePin.RIGHT);
  return [...left, ...middle, ...right].map((entry) => {
    const indexLeft = left.indexOf(entry);
    const indexRight = right.indexOf(entry);
    if (indexLeft >= 0)
      return {
        ...entry,
        offset: left.slice(0, indexLeft).reduce((sum, item) => sum + item.width, 0),
      };
    if (indexRight >= 0)
      return {
        ...entry,
        offset: right.slice(indexRight + 1).reduce((sum, item) => sum + item.width, 0),
      };
    return { ...entry, offset: 0 };
  });
};
