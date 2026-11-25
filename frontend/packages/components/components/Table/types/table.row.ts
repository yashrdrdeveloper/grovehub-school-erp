import type { ColumnPin, TableColumn } from './table.column';

export type TableEntry<T> =
  | { readonly kind: 'ROW'; readonly id: string; readonly row: T; readonly depth: number; readonly expandable: boolean }
  | { readonly kind: 'GROUP'; readonly id: string; readonly label: string; readonly count: number; readonly expanded: boolean }
  | { readonly kind: 'DETAIL'; readonly id: string; readonly row: T };

export interface TableColumnLayout<T> {
  readonly column: TableColumn<T>;
  readonly width: number;
  readonly pin: ColumnPin;
  readonly offset: number;
}
