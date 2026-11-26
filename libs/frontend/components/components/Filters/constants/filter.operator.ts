export const FilterOperator = Object.freeze({
  EQUALS: 'equals',
  NOT_EQUALS: 'notEquals',
  IN: 'in',
  NOT_IN: 'notIn',
  EXISTS: 'exists',
  NOT_EXISTS: 'notExists',
  CONTAINS: 'contains',
  STARTS_WITH: 'startsWith',
  ENDS_WITH: 'endsWith',
  GREATER_THAN: 'greaterThan',
  LESS_THAN: 'lessThan',
  BETWEEN: 'between',
} as const);

export type FilterOperator = (typeof FilterOperator)[keyof typeof FilterOperator];
