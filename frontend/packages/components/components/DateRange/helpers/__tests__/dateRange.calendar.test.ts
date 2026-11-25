import { DATE_RANGE_PRESETS } from '../../constants/dateRange.presets';
import { formatRangeDate, getRangeSelection, shiftRangeMonth } from '../dateRange.calendar';

describe('getRangeSelection', () => {
  it('selects the first date, orders reversed endpoints, and starts a new range after completion', () => {
    const start = getRangeSelection({ startDate: null, endDate: null }, '2026-11-20');
    expect(start).toEqual({ startDate: '2026-11-20', endDate: null });
    const complete = getRangeSelection(start, '2026-11-15');
    expect(complete).toEqual({ startDate: '2026-11-15', endDate: '2026-11-20' });
    expect(getRangeSelection(complete, '2026-11-25')).toEqual({ startDate: '2026-11-25', endDate: null });
    expect(getRangeSelection(start, '2026-11-20')).toEqual({ startDate: '2026-11-20', endDate: '2026-11-20' });
  });
});

describe('shiftRangeMonth', () => {
  it('keeps adjacent calendars synchronized across year boundaries', () => {
    expect(formatRangeDate(shiftRangeMonth(new Date(2026, 11, 31), 1))).toBe('2027-01-01');
  });
});

describe('DATE_RANGE_PRESETS', () => {
  it.each([
    ['today', '2026-11-26', '2026-11-26'],
    ['yesterday', '2026-11-25', '2026-11-25'],
    ['thisWeek', '2026-11-23', '2026-11-29'],
    ['lastWeek', '2026-11-16', '2026-11-22'],
    ['last7Days', '2026-11-20', '2026-11-26'],
    ['last30Days', '2026-10-28', '2026-11-26'],
    ['thisMonth', '2026-11-01', '2026-11-30'],
    ['lastMonth', '2026-10-01', '2026-10-31'],
    ['monthToDate', '2026-11-01', '2026-11-26'],
    ['quarterToDate', '2026-10-01', '2026-11-26'],
    ['yearToDate', '2026-01-01', '2026-11-26'],
  ])('computes %s using inclusive local calendar dates', (id, start, end) => {
    const preset = DATE_RANGE_PRESETS.find((item) => item.id === id)!;
    const today = new Date(2026, 10, 26, 12);
    expect(formatRangeDate(preset.start(today))).toBe(start);
    expect(formatRangeDate(preset.end(today))).toBe(end);
  });
});
