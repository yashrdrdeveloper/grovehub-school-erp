import { ResolvedTheme, ThemePreference } from '../types/useWorkspaceTheme';

export const WORKSPACE_THEME_BY_PREFERENCE: Readonly<
  Record<ThemePreference, ResolvedTheme | null>
> = {
  [ThemePreference.SYSTEM]: null,
  [ThemePreference.LIGHT]: ResolvedTheme.LIGHT,
  [ThemePreference.DARK]: ResolvedTheme.DARK,
};

export const WORKSPACE_THEME_BY_SYSTEM: Readonly<Record<string, ResolvedTheme>> = {
  true: ResolvedTheme.DARK,
  false: ResolvedTheme.LIGHT,
};

export const WORKSPACE_THEME_FROM_STORAGE: Readonly<Record<string, ThemePreference>> = {
  [ThemePreference.SYSTEM]: ThemePreference.SYSTEM,
  [ThemePreference.LIGHT]: ThemePreference.LIGHT,
  [ThemePreference.DARK]: ThemePreference.DARK,
};
