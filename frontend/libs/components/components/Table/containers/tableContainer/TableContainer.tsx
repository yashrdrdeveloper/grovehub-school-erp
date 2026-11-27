import { TableView } from '../../Table';
import { useTable } from '../../hooks/useTable';

import type { TableProps } from '../../types/table.props';

export const TableContainer = <T,>(props: TableProps<T>) => {
  const controls = useTable(props);
  return <TableView {...props} controls={controls} />;
};
