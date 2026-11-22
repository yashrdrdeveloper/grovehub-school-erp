import type { GetDashboardDateParts } from './types/dashboard.getDateParts';

const getGreeting = (hour: number): 'morning' | 'afternoon' | 'evening' => {
  if (hour < 12) return 'morning';
  if (hour < 17) return 'afternoon';
  return 'evening';
};

/** Pure calendar presentation helper; only date input affects output. */
export const getDashboardDateParts: GetDashboardDateParts = date => ({
  greeting: getGreeting(date.getHours()),
  dayName: new Intl.DateTimeFormat(undefined, { weekday: 'long' }).format(date),
  dateLabel: new Intl.DateTimeFormat(undefined, { day: 'numeric', month: 'long', year: 'numeric' }).format(date),
  dayOfMonth: date.getDate(),
});
