import lodashToLower from 'lodash-es/toLower';

import type { ToLower } from './types/toLower';

/** Typed case-and-word normalization through the common Base Lodash boundary. */
export const _toLower: ToLower = (value) => lodashToLower(value);
