import { FilterOperator } from '../constants/filter.operator';
import type { FilterCondition, FilterField } from '../types/filters.props';

export const isFilterValid = (condition: FilterCondition, field: FilterField | undefined): boolean => {
  if (!field || !field.operators.includes(condition.operator)) return false;
  if (condition.operator === FilterOperator.EXISTS || condition.operator === FilterOperator.NOT_EXISTS) {
    return condition.values.length === 0;
  }
  if (condition.operator === FilterOperator.BETWEEN) {
    return condition.values.length === 2 && condition.values.every((value) => value.trim().length > 0);
  }
  if (condition.operator === FilterOperator.IN || condition.operator === FilterOperator.NOT_IN) {
    return condition.values.length > 0 && condition.values.every((value) => value.trim().length > 0);
  }
  return condition.values.length === 1 && Boolean(condition.values[0]?.trim());
};
