import lodashCompact from 'lodash-es/compact';

import type { Compact } from './types/compact';

/** Typed Lodash compact gateway: excludes empty string, false, null and undefined. */
export const _compact: Compact = <T>(
  values: readonly (T | null | undefined | false | '')[],
): T[] => lodashCompact(values as (T | null | undefined | false | '')[]) as T[];
