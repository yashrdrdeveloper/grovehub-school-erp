import type { useTable } from '../hooks/useTable';
import type { TableProps } from './table.props';

export type TableViewProps<T> = TableProps<T> & {
  readonly controls: ReturnType<typeof useTable<T>>;
};
