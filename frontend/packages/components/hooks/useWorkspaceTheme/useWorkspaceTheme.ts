import { useEffect, useState } from 'react';

import { resolveWorkspaceTheme } from './helpers/workspaceTheme.resolve';
import type { ThemePreference, UseWorkspaceTheme } from './types/useWorkspaceTheme';

const STORAGE_KEY = 'groovehub.theme';
const PHONE_QUERY = '(max-width: 720px)';
const DARK_QUERY = '(prefers-color-scheme: dark)';

const readPreference = (): ThemePreference => {
  try {
    const value = window.localStorage.getItem(STORAGE_KEY);
    if (value === 'light' || value === 'dark') return value;
    return 'system';
  } catch {
    return 'system';
  }
};

export const useWorkspaceTheme: UseWorkspaceTheme = () => {
  const [preference, setPreferenceState] = useState<ThemePreference>(readPreference);

  useEffect(() => {
    const dark = window.matchMedia(DARK_QUERY);
    const phone = window.matchMedia(PHONE_QUERY);
    const apply = () => {
      const theme = resolveWorkspaceTheme({
        preference, systemDark: dark.matches, isMobile: phone.matches,
      });
      document.documentElement.dataset.theme = theme;
      document.documentElement.style.colorScheme = theme;
    };
    apply();
    dark.addEventListener('change', apply);
    phone.addEventListener('change', apply);
    return () => {
      dark.removeEventListener('change', apply);
      phone.removeEventListener('change', apply);
    };
  }, [preference]);

  const setPreference = (next: ThemePreference) => {
    try { window.localStorage.setItem(STORAGE_KEY, next); }
    catch { /* Browser storage is optional. */ }
    setPreferenceState(next);
  };
  return { preference, setPreference };
};
