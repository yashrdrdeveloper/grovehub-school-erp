import { TableDensity, TablePin } from '../constants/table.type';
import type { TablePreferences, TableProps } from '../types/table.props';
import type { TableColumn, ColumnPin } from '../types/table.column';

export const getDefaultTablePreferences = <T,>(props: TableProps<T>): TablePreferences => ({
  hidden: props.columns.filter((column) => column.visible === false && column.hideable !== false).map((column) => column.id),
  order: props.columns.map((column) => column.id),
  widths: {},
  pins: {},
  density: props.density ?? TableDensity.DEFAULT,
});

export const getTableColumnWidth = <T,>(column: TableColumn<T>, width?: number): number => {
  const fallback = column.width ?? 180;
  const requested = width ?? fallback;
  if (!Number.isFinite(requested)) return fallback;
  return Math.max(column.minWidth ?? 80, Math.min(column.maxWidth ?? 1200, requested));
};

export const readTablePreferences = <T,>(props: TableProps<T>): TablePreferences => {
  const defaults = getDefaultTablePreferences(props);
  const key = props.preferencesProps?.storageKey;
  if (!key || typeof window === 'undefined') return defaults;
  try {
    const raw = window.localStorage.getItem(key);
    if (!raw) return defaults;
    const saved = JSON.parse(raw) as Partial<TablePreferences>;
    const ids = defaults.order;
    const hidden = (Array.isArray(saved.hidden) && saved.hidden.filter((id) => ids.includes(id))) || defaults.hidden;
    const order = (Array.isArray(saved.order) && saved.order.filter((id) => ids.includes(id))) || ids;
    const widths: Record<string, number> = {};
    const pins: Record<string, ColumnPin> = {};
    for (const column of props.columns) {
      const width = saved.widths?.[column.id];
      const pin = saved.pins?.[column.id];
      if (typeof width === 'number') widths[column.id] = getTableColumnWidth(column, width);
      if (pin && Object.values(TablePin).includes(pin)) pins[column.id] = pin;
    }
    const density = saved.density;
    return {
      hidden: hidden.filter((id) => props.columns.find((column) => column.id === id)?.hideable !== false),
      order: Array.from(new Set([...order, ...ids])),
      widths,
      pins,
      density: Object.values(TableDensity).find((item) => item === density) ?? defaults.density,
    };
  } catch {
    return defaults;
  }
};
