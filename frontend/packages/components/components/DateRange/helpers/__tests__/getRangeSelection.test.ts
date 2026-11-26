import { getRangeSelection } from '../dateRange.calendar';

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
