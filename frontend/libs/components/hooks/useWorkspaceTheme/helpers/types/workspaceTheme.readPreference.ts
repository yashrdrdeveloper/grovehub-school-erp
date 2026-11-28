import type { ThemePreference } from '../../types/useWorkspaceTheme';

export type ReadThemePreference = (storageKey: string) => ThemePreference;
