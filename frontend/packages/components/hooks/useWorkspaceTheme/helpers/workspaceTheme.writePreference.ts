import type { WriteThemePreference } from './types/workspaceTheme.writePreference';

/** Local storage is best-effort. State updates remain available when storage is blocked. */
export const writeThemePreference: WriteThemePreference = (preference, storageKey) => {
  try {
    window.localStorage.setItem(storageKey, preference);
  } catch {
    // Browsers may disable local storage in private or restricted contexts.
  }
};
