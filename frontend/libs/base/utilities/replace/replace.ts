import lodashReplace from 'lodash-es/replace';

import type { Replace } from './types/replace';

/** Pure typed Lodash regular-expression replacement for Base consumers. */
export const _replace: Replace = (value, pattern, replacement) =>
  lodashReplace(value, pattern, replacement);
