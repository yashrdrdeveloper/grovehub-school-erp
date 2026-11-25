import { shiftRangeDay } from '../helpers/dateRange.calendar';
import type { DateRangePreset } from '../types/dateRange.preset';

const today = (date: Date) => date;
const weekStart = (date: Date) => shiftRangeDay(date, -((date.getDay() + 6) % 7));
const monthStart = (date: Date) => new Date(date.getFullYear(), date.getMonth(), 1);
const quarterStart = (date: Date) => new Date(date.getFullYear(), Math.floor(date.getMonth() / 3) * 3, 1);
const yearStart = (date: Date) => new Date(date.getFullYear(), 0, 1);

export const DATE_RANGE_PRESETS: readonly DateRangePreset[] = Object.freeze([
  { id: 'today', label: __('Today'), start: today, end: today },
  { id: 'yesterday', label: __('Yesterday'), start: (date) => shiftRangeDay(date, -1), end: (date) => shiftRangeDay(date, -1) },
  { id: 'thisWeek', label: __('This week'), start: weekStart, end: (date) => shiftRangeDay(weekStart(date), 6) },
  { id: 'lastWeek', label: __('Last week'), start: (date) => shiftRangeDay(weekStart(date), -7), end: (date) => shiftRangeDay(weekStart(date), -1) },
  { id: 'last7Days', label: __('Last 7 days'), start: (date) => shiftRangeDay(date, -6), end: today },
  { id: 'last30Days', label: __('Last 30 days'), start: (date) => shiftRangeDay(date, -29), end: today },
  { id: 'thisMonth', label: __('This month'), start: monthStart, end: (date) => new Date(date.getFullYear(), date.getMonth() + 1, 0) },
  { id: 'lastMonth', label: __('Last month'), start: (date) => new Date(date.getFullYear(), date.getMonth() - 1, 1), end: (date) => new Date(date.getFullYear(), date.getMonth(), 0) },
  { id: 'monthToDate', label: __('Month to date'), start: monthStart, end: today },
  { id: 'quarterToDate', label: __('Quarter to date'), start: quarterStart, end: today },
  { id: 'yearToDate', label: __('Year to date'), start: yearStart, end: today },
]);
