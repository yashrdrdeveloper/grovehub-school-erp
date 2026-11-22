import type { GetDashboardCalendarCopy } from './types/dashboardAside.getCalendarCopy';

export const getDashboardCalendarCopy: GetDashboardCalendarCopy = (date, today) => {
  if (today)
    return {
      calendarDay: String(date.dayOfMonth),
      calendarHeading: date.dayName,
      calendarDetail: date.dateLabel,
      calendarEmptyTitle: 'Today’s schedule',
    };
  return {
    calendarDay: '→',
    calendarHeading: 'Next on your calendar',
    calendarDetail: 'Calendar preview',
    calendarEmptyTitle: 'Upcoming events',
  };
};
