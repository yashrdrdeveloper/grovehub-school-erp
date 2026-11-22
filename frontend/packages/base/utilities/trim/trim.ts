import lodashTrim from 'lodash-es/trim';
import type { Trim } from './types/trim';

/** Pure typed gateway to Lodash's whitespace normalization. */
export const _trim: Trim = value => lodashTrim(value);
