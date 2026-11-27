import lodashFlatMap from 'lodash-es/flatMap';

import type { FlatMap } from './types/flatMap';

/** Typed order-preserving Lodash flatMap gateway for UI view models. */
export const _flatMap: FlatMap = <T, R>(
  values: readonly T[],
  iterator: (value: T, index: number) => readonly R[],
): R[] => lodashFlatMap(values as T[], iterator) as R[];
