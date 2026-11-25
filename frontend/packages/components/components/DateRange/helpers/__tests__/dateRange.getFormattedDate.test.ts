import { getFormattedRangeDate } from '../dateRange.getFormattedDate';

describe('getFormattedRangeDate', () => {
  it('formats calendar dates without time or shifting the selected day', () => {
    expect(getFormattedRangeDate('2026-11-01')).toBe('01 Nov 2026');
    expect(getFormattedRangeDate('2026-11-30')).toBe('30 Nov 2026');
    expect(getFormattedRangeDate('2024-02-29')).toBe('29 Feb 2024');
  });

  it('preserves placeholders for empty or invalid dates', () => {
    expect(getFormattedRangeDate(null, 'Start date')).toBe('Start date');
    expect(getFormattedRangeDate(undefined, 'End date')).toBe('End date');
    expect(getFormattedRangeDate('2026-02-30')).toBe('—');
  });
});
