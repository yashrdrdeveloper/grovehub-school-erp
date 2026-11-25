export const TableType = Object.freeze({
  DEFAULT: 'DEFAULT',
  VIRTUALIZED: 'VIRTUALIZED',
  EXPANDABLE: 'EXPANDABLE',
  GROUPED: 'GROUPED',
  TREE: 'TREE',
} as const);

export const TableDataMode = Object.freeze({
  SERVER: 'SERVER',
  CLIENT: 'CLIENT',
} as const);

export const TableDensity = Object.freeze({
  COMPACT: 'COMPACT',
  DEFAULT: 'DEFAULT',
  COMFORTABLE: 'COMFORTABLE',
} as const);

export const TablePin = Object.freeze({
  NONE: 'NONE',
  LEFT: 'LEFT',
  RIGHT: 'RIGHT',
} as const);

export const TableAlignment = Object.freeze({
  START: 'START',
  CENTER: 'CENTER',
  END: 'END',
} as const);
