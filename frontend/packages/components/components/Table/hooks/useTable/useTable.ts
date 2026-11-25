import { useState } from 'react';
import { _filter } from '@Base/utilities/filter';
import { _map } from '@Base/utilities/map';

import { TableDensity } from '../../constants/table.type';
import type { TableProps } from '../../types/table.props';

export const useTable = <T,>({ columns, density }: TableProps<T>) => {
  const [hidden, setHidden] = useState<readonly string[]>(() => _map(_filter(columns, (column) => column.visible === false), (column) => column.id));
  const [rowDensity, setRowDensity] = useState(density ?? TableDensity.DEFAULT);
  const visibleColumns = _filter(columns, (column) => !hidden.includes(column.id));
  const toggle = (id: string) => {
    const column = columns.find((entry) => entry.id === id);
    if (!column || column.hideable === false) return;
    if (hidden.includes(id)) {
      setHidden(_filter(hidden, (entry) => entry !== id));
      return;
    }
    if (visibleColumns.length <= 1) return;
    setHidden([...hidden, id]);
  };
  return {
    hidden, visibleColumns, rowDensity, setRowDensity, toggle,
    reset: () => {
      setHidden(_map(_filter(columns, (column) => column.visible === false), (column) => column.id));
      setRowDensity(density ?? TableDensity.DEFAULT);
    },
  };
};
