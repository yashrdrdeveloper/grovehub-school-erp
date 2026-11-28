import type { ReactNode } from 'react';
import type { TableAlignment, TablePin } from '../constants/table.type';

export type ColumnPin = (typeof TablePin)[keyof typeof TablePin];
export type ColumnAlignment = (typeof TableAlignment)[keyof typeof TableAlignment];

export interface TableColumn<T> {
  readonly id: string;
  readonly header: string;
  readonly accessor: (row: T) => unknown;
  readonly renderer?: (row: T) => ReactNode;
  readonly editRenderer?: (row: T, commit: (value: unknown) => void) => ReactNode;
  readonly compare?: (left: T, right: T) => number;
  readonly sortable?: boolean;
  readonly sortField?: string;
  readonly width?: number;
  /** Grow into spare viewport width until explicitly resized. */
  readonly fillRemainingWidth?: boolean;
  readonly minWidth?: number;
  readonly maxWidth?: number;
  readonly pin?: ColumnPin;
  readonly alignment?: ColumnAlignment;
  readonly hideable?: boolean;
  readonly visible?: boolean;
  readonly resizable?: boolean;
  readonly reorderable?: boolean;
  readonly mobileVisible?: boolean;
  readonly searchable?: boolean;
  readonly exportable?: boolean;
}
