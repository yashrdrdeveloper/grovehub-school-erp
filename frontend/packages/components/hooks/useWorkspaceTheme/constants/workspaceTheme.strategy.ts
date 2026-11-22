import { WORKSPACE_THEME_BY_PREFERENCE, WORKSPACE_THEME_BY_SYSTEM } from './workspaceTheme.resolution';
import { ThemeViewport } from './workspaceTheme.viewport';

import type { ResolvedTheme, ThemePreference } from '../types/useWorkspaceTheme';

type ThemeResolver = (
  preference: ThemePreference,
  systemDark: boolean,
) => ResolvedTheme;

const systemAppearance = (systemDark: boolean): ResolvedTheme =>
  WORKSPACE_THEME_BY_SYSTEM[String(systemDark)]!;

export const WORKSPACE_THEME_STRATEGIES: Readonly<Record<ThemeViewport, ThemeResolver>> = {
  [ThemeViewport.MOBILE]: (_preference, systemDark) => systemAppearance(systemDark),
  [ThemeViewport.DESKTOP]: (preference, systemDark) =>
    WORKSPACE_THEME_BY_PREFERENCE[preference] ?? systemAppearance(systemDark),
};
