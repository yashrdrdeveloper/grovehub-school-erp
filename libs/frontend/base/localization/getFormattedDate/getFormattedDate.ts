import type { FormattedDateOptions } from './types/getFormattedDate';

/** Format the supplied instant in the configured locale and timezone, without changing it. */
export const getFormattedDate = (value: string | number | Date | null | undefined, options: FormattedDateOptions = {}): string => {
  const fallback = options.fallback ?? '—';
  if (value === null || value === undefined || value === '') return fallback;
  const timestamp = () => {
    if (value instanceof Date) return value.getTime();
    return value;
  };
  const date = new Date(timestamp());
  if (!Number.isFinite(date.getTime())) return fallback;
  if (options.timeOnly) return new Intl.DateTimeFormat(options.locale ?? 'en-IN', { hour: '2-digit', minute: '2-digit', hour12: true, timeZone: options.timeZone }).format(date).toUpperCase();
  const format: Intl.DateTimeFormatOptions = { day: '2-digit', month: 'short', year: 'numeric', timeZone: options.timeZone };
  if (options.includeTime !== false) {
    format.hour = '2-digit';
    format.minute = '2-digit';
    format.hour12 = true;
  }
  return new Intl.DateTimeFormat(options.locale ?? 'en-IN', format).format(date);
};
