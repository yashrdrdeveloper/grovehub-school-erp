import type { TableViewProps } from '../../../types/table.view';

export interface TableRowActionsProps<T> {
  readonly row: T;
  readonly props: TableViewProps<T>;
}
