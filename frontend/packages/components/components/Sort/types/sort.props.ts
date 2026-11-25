import type { SortDirection } from '../constants/sort.direction';

export interface SortField {
  readonly id: string;
  readonly label: string;
}

export interface SortCondition {
  readonly field: string;
  readonly direction: SortDirection;
}

export interface SortProps {
  readonly fields: readonly SortField[];
  readonly value: readonly SortCondition[];
  readonly onChange: (value: readonly SortCondition[]) => void;
  readonly maxFields?: number;
  readonly iconOnly?: boolean;
  readonly disabled?: boolean;
  readonly dataTestId?: string;
}
