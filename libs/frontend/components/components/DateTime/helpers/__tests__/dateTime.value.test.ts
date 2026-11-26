import { getUTCDateTimeValue } from '../dateTime.value';

describe('getUTCDateTimeValue', () => {
  it('converts explicit offsets to canonical UTC and rejects invalid input', () => {
    expect(getUTCDateTimeValue('2026-11-25T10:30:00+05:30')).toBe('2026-11-25T05:00:00.000Z');
    expect(getUTCDateTimeValue('')).toBeNull();
    expect(getUTCDateTimeValue('invalid')).toBeNull();
  });
});
