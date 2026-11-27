export const getTableFillWidth = (available: number, columns: number, fixed: number): number =>
  Math.max(0, Math.floor(available) - columns - fixed);
