import { getAvatarHue } from '../avatar.getHue';

describe('getAvatarHue', () => {
  it('returns the same hue for normalized names', () => {
    expect(getAvatarHue(' Rahul Sharma ')).toBe(getAvatarHue('rahul sharma'));
  });

  it('keeps hues in range and varies by name', () => {
    const hues = ['Rahul Sharma', 'Priya Gupta', 'Unknown user'].map(getAvatarHue);
    expect(hues.every((hue) => hue >= 0 && hue < 360)).toBe(true);
    expect(new Set(hues).size).toBeGreaterThan(1);
  });
});
