import { useCallback, useEffect, useState } from 'react';

import {
  WORKSPACE_THEME_DARK_QUERY,
  WORKSPACE_THEME_PHONE_QUERY,
} from './constants/workspaceTheme.storage';
import { readThemePreference } from './helpers/workspaceTheme.readPreference';
import { resolveWorkspaceTheme } from './helpers/workspaceTheme.resolve';
import { writeThemePreference } from './helpers/workspaceTheme.writePreference';

import type { ThemePreference, UseWorkspaceTheme } from './types/useWorkspaceTheme';

/** One theme owner for persistent desktop preference and system-based mobile behavior. */
export const useWorkspaceTheme: UseWorkspaceTheme = () => {
  const [preference, setPreferenceState] = useState<ThemePreference>(readThemePreference);

  useEffect(() => {
    const dark = window.matchMedia(WORKSPACE_THEME_DARK_QUERY);
    const phone = window.matchMedia(WORKSPACE_THEME_PHONE_QUERY);

    const applyTheme = () => {
      const theme = resolveWorkspaceTheme({
        preference,
        systemDark: dark.matches,
        isMobile: phone.matches,
      });

      document.documentElement.dataset.theme = theme;
      document.documentElement.style.colorScheme = theme;
    };

    applyTheme();
    dark.addEventListener('change', applyTheme);
    phone.addEventListener('change', applyTheme);

    return () => {
      dark.removeEventListener('change', applyTheme);
      phone.removeEventListener('change', applyTheme);
    };
  }, [preference]);

  const setPreference = useCallback((next: ThemePreference) => {
    writeThemePreference(next);
    setPreferenceState(next);
  }, []);

  return {
    preference,
    setPreference,
  };
};
