import { getFormattedDate } from '../getFormattedDate';

describe('getFormattedDate', () => {
  it('formats the original instant in the requested timezone', () => {
    const value = getFormattedDate('2026-11-25T05:00:00Z', { locale: 'en-IN', timeZone: 'Asia/Kolkata' });
    expect(value).toContain('25 Nov 2026');
    expect(value).toMatch(/10:30\s*am/i);
  });

  it('does not replace invalid or missing timestamps with the current time', () => {
    expect(getFormattedDate('invalid')).toBe('—');
    expect(getFormattedDate(undefined)).toBe('—');
  });
});
