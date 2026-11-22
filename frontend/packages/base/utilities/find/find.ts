import lodashFind from 'lodash-es/find';
import type { Find } from './types/find';

/** Preserve readonly inputs while returning the first matching element. */
export const _find: Find = (values, predicate) => lodashFind(values as unknown[], predicate);
