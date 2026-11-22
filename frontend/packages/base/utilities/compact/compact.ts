import lodashCompact from 'lodash-es/compact';

import type { Compact } from './types/compact';

/** Pure, typed Base gateway to Lodash; feature packages import only @Base. */
export const _compact: Compact = (values) =>
  lodashCompact(values as (T | null | undefined | false | '')[]);
