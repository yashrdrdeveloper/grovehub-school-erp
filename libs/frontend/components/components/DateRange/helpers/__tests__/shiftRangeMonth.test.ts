import { formatRangeDate, shiftRangeMonth } from '../dateRange.calendar';

describe('shiftRangeMonth', () => {
  it('keeps adjacent calendars synchronized across year boundaries', () => {
    expect(formatRangeDate(shiftRangeMonth(new Date(2026, 11, 31), 1))).toBe('2027-01-01');
  });
});
