import lodashIsEqual from 'lodash-es/isEqual';

import type { IsEqual } from './types/isEqual';

export const isEqual: IsEqual = (left, right) => lodashIsEqual(left, right);
