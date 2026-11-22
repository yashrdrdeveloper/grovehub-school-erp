import type { GetThemeChoiceClassName } from './types/themeChoices.getClassName';
export const getThemeChoiceClassName: GetThemeChoiceClassName = (active) => {
  if (active) return 'theme-choice theme-choice--active';
  return 'theme-choice';
};
