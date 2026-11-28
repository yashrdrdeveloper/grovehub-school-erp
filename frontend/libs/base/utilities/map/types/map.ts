export type MapValues = <T, R>(
  values: readonly T[],
  iterator: (value: T, index: number) => R,
) => R[];
