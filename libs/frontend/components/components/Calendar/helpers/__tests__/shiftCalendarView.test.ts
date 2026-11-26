import { formatCalendarDate } from '../calendar.date';
import { shiftCalendarView } from '../calendar.navigation';

describe('shiftCalendarView', () => {
  it('navigates months, years and decades without retaining the day of month', () => {
    const date = new Date(2026, 11, 31);
    expect(formatCalendarDate(shiftCalendarView(date, 'day', 1))).toBe('2027-01-01');
    expect(formatCalendarDate(shiftCalendarView(date, 'month', -1))).toBe('2025-12-01');
    expect(formatCalendarDate(shiftCalendarView(date, 'year', 1))).toBe('2036-12-01');
  });
});
