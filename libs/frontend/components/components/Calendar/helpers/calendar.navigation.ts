import type { CalendarMode } from '../types/calendar.mode';

export const getCalendarYears = (year: number): readonly number[] => {
  const start = Math.floor(year / 10) * 10;
  return Array.from({ length: 12 }, (_, index) => start - 1 + index);
};

export const shiftCalendarView = (month: Date, mode: CalendarMode, direction: number): Date => {
  if (mode === 'year') return new Date(month.getFullYear() + direction * 10, month.getMonth(), 1);
  if (mode === 'month') return new Date(month.getFullYear() + direction, month.getMonth(), 1);
  return new Date(month.getFullYear(), month.getMonth() + direction, 1);
};
