import lodashMap from 'lodash-es/map';
import type { MapValues } from './types/map';

/** Pure, typed Lodash gateway. */
export const _map = <T, R>(values: readonly T[], iterator: (value: T, index: number) => R): R[] =>
  lodashMap(values as T[], iterator);

const _mapContract: MapValues = _map;
void _mapContract;
