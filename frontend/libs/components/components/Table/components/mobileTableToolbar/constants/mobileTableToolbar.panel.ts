export const MobileTablePanel = {
  OPTIONS: 'OPTIONS',
  FILTER: 'FILTER',
  SORT: 'SORT',
  SETTINGS: 'SETTINGS',
  BATCH: 'BATCH',
} as const;

export type MobileTablePanel = (typeof MobileTablePanel)[keyof typeof MobileTablePanel];
