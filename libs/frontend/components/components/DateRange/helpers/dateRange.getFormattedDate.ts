import { getFormattedDate } from '@Base/localization/getFormattedDate';

import { parseCalendarDate } from '../../Calendar/helpers/calendar.date';

export const getFormattedRangeDate = (value: string | null | undefined, fallback = '—'): string =>
  getFormattedDate(parseCalendarDate(value ?? null), { includeTime: false, fallback });
