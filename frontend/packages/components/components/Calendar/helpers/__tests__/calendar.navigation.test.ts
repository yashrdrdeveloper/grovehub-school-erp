import { getCalendarYears, shiftCalendarView } from '../calendar.navigation';
import { formatCalendarDate } from '../calendar.date';

describe('getCalendarYears', () => {
  it('returns a decade and the adjacent years', () => {
    expect(getCalendarYears(2026)).toEqual([2019, 2020, 2021, 2022, 2023, 2024, 2025, 2026, 2027, 2028, 2029, 2030]);
  });
});

describe('shiftCalendarView', () => {
  it('navigates months, years and decades without retaining the day of month', () => {
    const date = new Date(2026, 11, 31);
    expect(formatCalendarDate(shiftCalendarView(date, 'day', 1))).toBe('2027-01-01');
    expect(formatCalendarDate(shiftCalendarView(date, 'month', -1))).toBe('2025-12-01');
    expect(formatCalendarDate(shiftCalendarView(date, 'year', 1))).toBe('2036-12-01');
  });
});
