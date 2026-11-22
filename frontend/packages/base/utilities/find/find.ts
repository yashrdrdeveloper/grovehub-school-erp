import lodashFind from 'lodash-es/find';
import type { Find } from './types/find';

/** Pure, typed Lodash gateway. */
export const _find = <T>(values: readonly T[], predicate: (value: T) => boolean): T | undefined =>
  lodashFind(values as T[], predicate);

const _findContract: Find = _find;
void _findContract;
