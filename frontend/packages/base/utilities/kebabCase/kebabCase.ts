import lodashKebabCase from 'lodash-es/kebabCase';

import type { KebabCase } from './types/kebabCase';

/** Pure, typed Base gateway to Lodash; feature packages import only @Base. */
export const _kebabCase: KebabCase = (value) =>
  lodashKebabCase(value);
