import type { TableColumn } from '../../../types/table.column';
import type { TableActions } from '../../../types/table.props';

export interface TableCellProps<T> {
  readonly row: T;
  readonly column: TableColumn<T>;
  readonly actions?: TableActions<T>;
}
