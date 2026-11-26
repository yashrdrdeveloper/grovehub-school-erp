export const getFormattedNumber = (value: number, locale = 'en-IN', options: Intl.NumberFormatOptions = {}): string => {
  if (!Number.isFinite(value)) return '—';
  return new Intl.NumberFormat(locale, options).format(value);
};
