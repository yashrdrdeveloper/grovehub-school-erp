import type { ReactNode } from 'react';

export interface RecordColumn<T> {
  readonly id: string;
  readonly title: string;
  readonly cell: (record: T) => ReactNode;
}
export interface PagedRecordsProps<T> {
  readonly resource: string;
  readonly revision: number;
  readonly columns: readonly RecordColumn<T>[];
  readonly rowKey: (record: T) => string | number;
  readonly onSelect?: (record: T) => void;
  readonly title: string;
}
