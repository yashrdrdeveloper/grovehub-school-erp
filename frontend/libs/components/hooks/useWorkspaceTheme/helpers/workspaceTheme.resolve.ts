import { THEME_VIEWPORT_FROM_MEDIA } from '../constants/workspaceTheme.viewport';
import { WORKSPACE_THEME_STRATEGIES } from '../constants/workspaceTheme.strategy';

import type { ResolveWorkspaceTheme } from './types/workspaceTheme.resolve';

/** Enum-indexed strategies avoid conditional chains and preserve desktop preferences. */
export const resolveWorkspaceTheme: ResolveWorkspaceTheme = ({
  preference,
  systemDark,
  isMobile,
}) => {
  const viewport = THEME_VIEWPORT_FROM_MEDIA[String(isMobile)]!;
  return WORKSPACE_THEME_STRATEGIES[viewport](preference, systemDark);
};
