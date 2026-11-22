export type FlatMap = <T, R>(
  values: readonly T[],
  iterator: (value: T, index: number) => readonly R[],
) => R[];
