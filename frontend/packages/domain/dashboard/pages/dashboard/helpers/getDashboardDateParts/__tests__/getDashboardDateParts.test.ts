import { getDashboardDateParts } from '../getDashboardDateParts';

describe('getDashboardDateParts', () => {
  it('assigns a morning greeting before noon', () => {
    const result = getDashboardDateParts(new Date(2026, 10, 22, 9));
    expect(result.greeting).toBe('morning');
    expect(result.dayOfMonth).toBe(22);
    expect(result.dayName.length).toBeGreaterThan(0);
  });
  it('changes greeting at noon and in the evening', () => {
    expect(getDashboardDateParts(new Date(2026, 10, 22, 12)).greeting).toBe('afternoon');
    expect(getDashboardDateParts(new Date(2026, 10, 22, 17)).greeting).toBe('evening');
  });
});
