/** Closed, exported appearance options; feature code never invents string cases. */
export enum ThemePreference {
  SYSTEM = 'system',
  LIGHT = 'light',
  DARK = 'dark',
}

export enum ResolvedTheme {
  LIGHT = 'light',
  DARK = 'dark',
}

export interface WorkspaceThemeController {
  readonly preference: ThemePreference;
  readonly setPreference: (next: ThemePreference) => void;
}

/** The storage key is fixed for the lifetime of the mounted application. */
export type UseWorkspaceTheme = (storageKey: string) => WorkspaceThemeController;
