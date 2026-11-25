import { formatCalendarDate, parseCalendarDate } from '../../Calendar/helpers/calendar.date';

import type { DateRangeValue } from '../types/dateRange.props';

export const formatRangeDate = formatCalendarDate;

export const shiftRangeDay = (date: Date, offset: number): Date =>
  new Date(date.getFullYear(), date.getMonth(), date.getDate() + offset);

export const shiftRangeMonth = (date: Date, offset: number): Date =>
  new Date(date.getFullYear(), date.getMonth() + offset, 1);

export const getRangeMonth = (date: string | null, fallback: Date): Date => {
  return shiftRangeMonth(parseCalendarDate(date) ?? fallback, 0);
};

export const getRangeSelection = (range: DateRangeValue, date: string): DateRangeValue => {
  if (!range.startDate || range.endDate) return { startDate: date, endDate: null };
  if (date < range.startDate) return { startDate: date, endDate: range.startDate };
  return { startDate: range.startDate, endDate: date };
};
