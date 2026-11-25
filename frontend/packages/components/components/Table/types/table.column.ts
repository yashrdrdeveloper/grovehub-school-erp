import type { ReactNode } from 'react';

export interface TableColumn<T> {
  readonly id: string;
  readonly header: string;
  readonly accessor: (row: T) => unknown;
  readonly renderer?: (row: T) => ReactNode;
  readonly sortable?: boolean;
  readonly sortField?: string;
  readonly width?: number;
  readonly hideable?: boolean;
  readonly visible?: boolean;
}
