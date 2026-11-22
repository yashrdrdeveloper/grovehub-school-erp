export type ThemePreference = 'system' | 'light' | 'dark';
export type ResolvedTheme = 'light' | 'dark';

export interface WorkspaceThemeController {
  readonly preference: ThemePreference;
  readonly setPreference: (next: ThemePreference) => void;
}

export type UseWorkspaceTheme = () => WorkspaceThemeController;
