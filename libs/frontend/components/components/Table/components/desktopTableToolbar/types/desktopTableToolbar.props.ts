import type { ColumnManagerProps } from '../../../../ColumnManager';
import type { TableViewProps } from '../../../types/table.view';

export type DesktopTableToolbarProps<T> = TableViewProps<T> & {
  readonly options: ColumnManagerProps['options'];
};
