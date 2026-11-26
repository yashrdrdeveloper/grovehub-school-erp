import type { ThemePreference } from '../../types/useWorkspaceTheme';

export type WriteThemePreference = (preference: ThemePreference, storageKey?: string) => void;
