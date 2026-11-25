import type { TableColumn } from '../../../types/table.column';
import type { TableViewProps } from '../../../types/table.view';

export interface TableSortHeaderProps<T> {
  readonly column: TableColumn<T>;
  readonly props: TableViewProps<T>;
}
