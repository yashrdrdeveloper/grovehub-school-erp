import { getTableColumnWidth } from './table.preferences';
import type { TableViewProps } from '../types/table.view';

export const getTableColumnOptions = <T>(props: TableViewProps<T>) => {
  const { controls, columns } = props;
  const order = Array.from(
    new Set([...controls.preferences.order, ...columns.map((column) => column.id)]),
  );
  const options = order.flatMap((id) => {
    const column = columns.find((entry) => entry.id === id);
    if (!column) return [];
    return [
      {
        id,
        label: column.header,
        visible: controls.layout.some((entry) => entry.column.id === id),
        disabled: column.hideable === false,
        width: getTableColumnWidth(column, controls.preferences.widths[id]),
        pin: controls.preferences.pins[id] ?? column.pin,
        resizable: column.resizable,
        reorderable: column.reorderable,
      },
    ];
  });
  return options;
};
