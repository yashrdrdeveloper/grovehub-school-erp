import type { GetDashboardDateParts } from './types/getDashboardDateParts';

/** Pure calendar presentation helper; only date input affects output. */
export const getDashboardDateParts: GetDashboardDateParts = date => ({
  greeting: date.getHours() < 12 ? 'morning' : date.getHours() < 17 ? 'afternoon' : 'evening',
  dayName: new Intl.DateTimeFormat(undefined, { weekday: 'long' }).format(date),
  dateLabel: new Intl.DateTimeFormat(undefined, { day: 'numeric', month: 'long', year: 'numeric' }).format(date),
  dayOfMonth: date.getDate(),
});
