export const SortDirection = Object.freeze({
  ASC: 'ASC',
  DESC: 'DESC',
} as const);

export type SortDirection = (typeof SortDirection)[keyof typeof SortDirection];
