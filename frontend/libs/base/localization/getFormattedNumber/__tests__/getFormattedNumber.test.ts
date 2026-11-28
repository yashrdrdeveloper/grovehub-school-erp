import { getFormattedNumber } from '../getFormattedNumber';

describe('getFormattedNumber', () => {
  it('formats localized grouping', () => {
    expect(getFormattedNumber(13323)).toBe('13,323');
    expect(getFormattedNumber(1234567)).toBe('12,34,567');
  });
});
