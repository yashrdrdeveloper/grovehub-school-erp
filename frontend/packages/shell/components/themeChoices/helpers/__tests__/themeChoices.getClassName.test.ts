import { getThemeChoiceClassName } from '../themeChoices.getClassName';

describe('getThemeChoiceClassName', () => {
  it('returns the active modifier only for the selected option', () => {
    expect(getThemeChoiceClassName(true)).toContain('theme-choice--active');
    expect(getThemeChoiceClassName(false)).toBe('theme-choice');
  });
});
