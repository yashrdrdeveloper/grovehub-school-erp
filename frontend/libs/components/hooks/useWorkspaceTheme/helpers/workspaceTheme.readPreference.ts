import { WORKSPACE_THEME_FROM_STORAGE } from '../constants/workspaceTheme.resolution';
import { ThemePreference } from '../types/useWorkspaceTheme';

import type { ReadThemePreference } from './types/workspaceTheme.readPreference';

/** Reading an unavailable or unrecognized local preference always falls back to system. */
export const readThemePreference: ReadThemePreference = (storageKey) => {
  try {
    const saved = window.localStorage.getItem(storageKey);
    return WORKSPACE_THEME_FROM_STORAGE[saved ?? ''] ?? ThemePreference.SYSTEM;
  } catch {
    return ThemePreference.SYSTEM;
  }
};
