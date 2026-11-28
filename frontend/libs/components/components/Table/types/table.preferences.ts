import type { TableDensity } from '../constants/table.type';
import type { ColumnPin } from './table.column';

export type Density = (typeof TableDensity)[keyof typeof TableDensity];

export interface TablePreferences {
  readonly hidden: readonly string[];
  readonly order: readonly string[];
  readonly widths: Readonly<Record<string, number>>;
  readonly pins: Readonly<Record<string, ColumnPin>>;
  readonly density: Density;
}

export interface TablePreferencesProps {
  /** Caller must scope by tenant and account. No business data is persisted. */
  readonly storageKey?: string;
  readonly value?: TablePreferences;
  readonly onChange?: (value: TablePreferences) => void;
}
