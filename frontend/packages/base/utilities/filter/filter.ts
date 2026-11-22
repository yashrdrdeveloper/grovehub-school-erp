import lodashFilter from 'lodash-es/filter';
import type { Filter } from './types/filter';

export const _filter = <T>(values: readonly T[], predicate: (item: T) => boolean): T[] =>
  lodashFilter(values as T[], predicate);
const _filterContract: Filter = _filter;
void _filterContract;
