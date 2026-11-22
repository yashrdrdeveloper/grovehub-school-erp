import lodashJoin from 'lodash-es/join';

import type { Join } from './types/join';

/** Pure, typed Base gateway to Lodash; feature packages import only @Base. */
export const _join: Join = (values, separator) => lodashJoin(values as string[], separator);
