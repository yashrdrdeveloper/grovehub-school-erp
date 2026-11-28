import { getDashboardCalendarCopy } from '../dashboardAside.getCalendarCopy';

describe('getDashboardCalendarCopy', () => {
  const date = {
    dayOfMonth: 23,
    dayName: 'Monday',
    dateLabel: '23 November',
    greeting: 'morning' as const,
  };
  it('uses actual date content for the today tab', () => {
    expect(getDashboardCalendarCopy(date, true)).toMatchObject({
      calendarDay: '23',
      calendarHeading: 'Monday',
    });
  });
  it('uses non-fabricated placeholders for future events', () => {
    expect(getDashboardCalendarCopy(date, false)).toMatchObject({
      calendarDay: '→',
      calendarEmptyTitle: 'Upcoming events',
    });
  });
});
