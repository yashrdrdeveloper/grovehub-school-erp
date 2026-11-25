import type { FilterCondition, FilterField } from '../../../types/filters.props';

export interface FilterValueProps {
  readonly hideLabel?: boolean;
  readonly field: FilterField;
  readonly condition: FilterCondition;
  readonly onChange: (values: readonly string[]) => void;
}
