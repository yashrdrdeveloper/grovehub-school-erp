import type { ReactElement } from 'react';
import type { ColumnDef } from '@tanstack/react-table';

export interface DataTableProps<T> {
  readonly rows: T[];
  readonly columns: ColumnDef<T, unknown>[];
  readonly emptyTitle?: string;
  readonly dataTestId?: string;
}
export type DataTableComponent = <T>(props: DataTableProps<T>) => ReactElement;
