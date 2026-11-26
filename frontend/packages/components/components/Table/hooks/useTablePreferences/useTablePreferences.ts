import { useEffect, useState } from 'react';
import { _filter } from '@Base/utilities/filter';
import { _map } from '@Base/utilities/map';

import { TablePin } from '../../constants/table.type';
import { getTableColumnLayout } from '../../helpers/table.getColumnLayout';
import {
  getDefaultTablePreferences,
  getTableColumnWidth,
  readTablePreferences,
} from '../../helpers/table.preferences';
import type { ColumnPin } from '../../types/table.column';
import type { Density, TablePreferences, TableProps } from '../../types/table.props';

export const useTablePreferences = <T>(props: TableProps<T>) => {
  const { columns, actions, preferencesProps } = props;
  const [localPreferences, setLocalPreferences] = useState(() => readTablePreferences(props));
  const preferences = preferencesProps?.value ?? localPreferences;
  const layout = getTableColumnLayout(columns, preferences);
  const update = (next: TablePreferences) => {
    setLocalPreferences(next);
    preferencesProps?.onChange?.(next);
    if (preferencesProps?.storageKey) {
      try {
        window.localStorage.setItem(preferencesProps.storageKey, JSON.stringify(next));
      } catch {
        // Storage is optional; the current table remains usable.
      }
    }
  };
  useEffect(() => {
    setLocalPreferences(readTablePreferences(props));
  }, [preferencesProps?.storageKey, props.id]);
  const toggle = (id: string) => {
    const column = columns.find((entry) => entry.id === id);
    if (!column || column.hideable === false) return;
    const hidden = _filter(preferences.hidden, (entry) => entry !== id);
    if (!preferences.hidden.includes(id)) {
      if (layout.length <= 1) return;
      hidden.push(id);
    }
    update({ ...preferences, hidden });
    actions?.COLUMN_VISIBILITY?.(hidden);
  };
  const reorder = (ids: readonly string[]) => {
    const order = Array.from(new Set(ids));
    if (order.length !== columns.length || columns.some((column) => !order.includes(column.id)))
      return;
    const previous = Array.from(
      new Set([...preferences.order, ..._map(columns, (column) => column.id)]),
    );
    if (
      columns.some(
        (column) =>
          column.reorderable === false && previous.indexOf(column.id) !== order.indexOf(column.id),
      )
    )
      return;
    update({ ...preferences, order });
    actions?.COLUMN_REORDER?.(order);
  };
  const move = (id: string, delta: number) => {
    const order = Array.from(
      new Set([...preferences.order, ..._map(columns, (column) => column.id)]),
    );
    const from = order.indexOf(id);
    const target = from + delta;
    const other = order[target];
    if (from < 0 || !other) return;
    if (columns.find((column) => column.id === id)?.reorderable === false) return;
    if (columns.find((column) => column.id === other)?.reorderable === false) return;
    order[from] = other;
    order[target] = id;
    update({ ...preferences, order });
    actions?.COLUMN_REORDER?.(order);
  };
  const resize = (id: string, requested: number) => {
    const column = columns.find((entry) => entry.id === id);
    if (!column || column.resizable === false) return;
    const width = getTableColumnWidth(column, requested);
    update({ ...preferences, widths: { ...preferences.widths, [id]: width } });
    actions?.COLUMN_RESIZE?.({ id, width });
  };
  const pin = (id: string, value: ColumnPin) => {
    if (!columns.some((column) => column.id === id) || !Object.values(TablePin).includes(value))
      return;
    const locked = _filter(
      columns,
      (column) => (preferences.pins[column.id] ?? column.pin ?? TablePin.NONE) !== TablePin.NONE,
    );
    if (
      value !== TablePin.NONE &&
      !locked.some((column) => column.id === id) &&
      locked.length >= (props.maxPinnedColumns ?? 3)
    )
      return;
    update({ ...preferences, pins: { ...preferences.pins, [id]: value } });
    actions?.COLUMN_PIN?.({ id, pin: value });
  };
  const setRowDensity = (density: Density) => {
    update({ ...preferences, density });
    actions?.DENSITY_CHANGE?.(density);
  };
  return {
    preferences,
    layout,
    toggle,
    move,
    reorder,
    resize,
    pin,
    setRowDensity,
    reset: () => {
      update(getDefaultTablePreferences(props));
      actions?.RESET?.();
    },
  };
};
