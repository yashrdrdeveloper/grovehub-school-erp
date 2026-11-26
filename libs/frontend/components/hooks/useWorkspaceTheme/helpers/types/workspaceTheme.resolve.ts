import type { ThemePreference, ResolvedTheme } from '../../types/useWorkspaceTheme';

export interface ResolveWorkspaceThemeParams {
  readonly preference: ThemePreference;
  readonly systemDark: boolean;
  readonly isMobile: boolean;
}
export type ResolveWorkspaceTheme = (params: ResolveWorkspaceThemeParams) => ResolvedTheme;
