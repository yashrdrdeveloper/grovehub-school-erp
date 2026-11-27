import lodashFilter from 'lodash-es/filter';

import type { Filter } from './types/filter';

/** Preserve user-defined type predicates while filtering through Lodash. */
export function _filter<T, S extends T>(
  values: readonly T[],
  predicate: (item: T) => item is S,
): S[];

export function _filter<T>(values: readonly T[], predicate: (item: T) => boolean): T[];

export function _filter<T>(values: readonly T[], predicate: (item: T) => boolean): T[] {
  return lodashFilter(values as T[], predicate);
}

const _filterContract: Filter = _filter;
void _filterContract;
