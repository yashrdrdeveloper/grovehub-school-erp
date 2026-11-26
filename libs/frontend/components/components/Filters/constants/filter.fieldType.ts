export const FilterFieldType = Object.freeze({
  TEXT: 'TEXT',
  NUMBER: 'NUMBER',
  BOOLEAN: 'BOOLEAN',
  SELECT: 'SELECT',
  MULTI_SELECT: 'MULTI_SELECT',
  ASYNC_SELECT: 'ASYNC_SELECT',
  ASYNC_MULTI_SELECT: 'ASYNC_MULTI_SELECT',
  DATE: 'DATE',
  DATE_TIME: 'DATE_TIME',
  DATE_RANGE: 'DATE_RANGE',
  NUMBER_RANGE: 'NUMBER_RANGE',
  CUSTOM: 'CUSTOM',
} as const);

export type FilterFieldType = (typeof FilterFieldType)[keyof typeof FilterFieldType];
