import { DATE_RANGE_PRESETS } from '../../constants/dateRange.presets';
import { formatRangeDate } from '../dateRange.calendar';

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
