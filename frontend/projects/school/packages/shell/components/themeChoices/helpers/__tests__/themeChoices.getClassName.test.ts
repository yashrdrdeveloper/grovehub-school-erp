import { getThemeChoiceClassName } from '../themeChoices.getClassName';

describe('getThemeChoiceClassName', () => {
  it('only adds the owned selected state for active preference', () => {
    expect(getThemeChoiceClassName(true)).toContain('item');
    expect(getThemeChoiceClassName(true)).toContain('active');
    expect(getThemeChoiceClassName(false)).toBe('item');
  });
});
