import { parseCalendarDate } from '../calendar.date';

describe('parseCalendarDate', () => {
  it('rejects impossible calendar dates', () => {
    expect(parseCalendarDate('2026-02-29')).toBeNull();
    expect(parseCalendarDate('2026-04-31')).toBeNull();
    expect(parseCalendarDate('2024-02-29')?.getDate()).toBe(29);
  });
});
