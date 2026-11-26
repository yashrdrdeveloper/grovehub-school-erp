import { getCalendarYears } from '../calendar.navigation';

describe('getCalendarYears', () => {
  it('returns a decade and the adjacent years', () => {
    expect(getCalendarYears(2026)).toEqual([2019, 2020, 2021, 2022, 2023, 2024, 2025, 2026, 2027, 2028, 2029, 2030]);
  });
});
