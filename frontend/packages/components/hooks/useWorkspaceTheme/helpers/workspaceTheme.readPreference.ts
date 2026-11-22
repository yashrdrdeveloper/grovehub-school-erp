import { WORKSPACE_THEME_FROM_STORAGE } from '../constants/workspaceTheme.resolution';
import { WORKSPACE_THEME_STORAGE_KEY } from '../constants/workspaceTheme.storage';
import { ThemePreference } from '../types/useWorkspaceTheme';

import type { ReadThemePreference } from './types/workspaceTheme.readPreference';

/** Reading an unavailable or unrecognized local preference always falls back to system. */
export const readThemePreference: ReadThemePreference = () => {
  try {
    const saved = window.localStorage.getItem(WORKSPACE_THEME_STORAGE_KEY);
    return WORKSPACE_THEME_FROM_STORAGE[saved ?? ''] ?? ThemePreference.SYSTEM;
  } catch {
    return ThemePreference.SYSTEM;
  }
};
