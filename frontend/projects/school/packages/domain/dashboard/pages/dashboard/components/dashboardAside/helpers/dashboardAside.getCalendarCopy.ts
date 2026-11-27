import type { GetDashboardCalendarCopy } from './types/dashboardAside.getCalendarCopy';

export const getDashboardCalendarCopy: GetDashboardCalendarCopy = (date, today) => {
  if (today)
    return {
      calendarDay: String(date.dayOfMonth),
      calendarHeading: date.dayName,
      calendarDetail: date.dateLabel,
      calendarEmptyTitle: __('Today’s schedule'),
    };
  return {
    calendarDay: '→',
    calendarHeading: __('Next on your calendar'),
    calendarDetail: __('Calendar preview'),
    calendarEmptyTitle: __('Upcoming events'),
  };
};
