import { FilterOperator } from './filter.operator';

export const FILTER_OPERATOR_LABELS: Readonly<Record<FilterOperator, string>> = Object.freeze({
  [FilterOperator.EQUALS]: __('Equals'),
  [FilterOperator.NOT_EQUALS]: __('Does not equal'),
  [FilterOperator.IN]: __('In'),
  [FilterOperator.NOT_IN]: __('Not in'),
  [FilterOperator.EXISTS]: __('Exists'),
  [FilterOperator.NOT_EXISTS]: __('Does not exist'),
  [FilterOperator.CONTAINS]: __('Contains'),
  [FilterOperator.STARTS_WITH]: __('Starts with'),
  [FilterOperator.ENDS_WITH]: __('Ends with'),
  [FilterOperator.GREATER_THAN]: __('Greater than'),
  [FilterOperator.LESS_THAN]: __('Less than'),
  [FilterOperator.BETWEEN]: __('Between'),
});
