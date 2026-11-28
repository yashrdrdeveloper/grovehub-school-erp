import type { ReactNode } from 'react';
import type { SelectOption } from '../../Select';
import type { AsyncSelectSource } from '../../AsyncSelect';
import type { FilterFieldType } from '../constants/filter.fieldType';
import type { FilterOperator } from '../constants/filter.operator';

export interface FilterCondition {
  readonly id: string;
  readonly field: string;
  readonly operator: FilterOperator;
  readonly values: readonly string[];
  readonly dateRangePresetId?: string;
}

export interface FilterField {
  readonly id: string;
  readonly label: string;
  readonly type: FilterFieldType;
  readonly operators: readonly FilterOperator[];
  readonly options?: readonly SelectOption[];
  readonly source?: AsyncSelectSource<string>;
  readonly renderValue?: (
    condition: FilterCondition,
    change: (values: readonly string[]) => void,
  ) => ReactNode;
}

export interface FiltersProps {
  readonly fields: readonly FilterField[];
  readonly value: readonly FilterCondition[];
  readonly onChange: (value: readonly FilterCondition[]) => void;
  readonly maxFilters?: number;
  readonly iconOnly?: boolean;
  readonly disabled?: boolean;
  readonly dataTestId?: string;
  /** Mount as an already-open mobile editor, with no local trigger. */
  readonly initialOpen?: boolean;
  readonly onDismiss?: () => void;
}
