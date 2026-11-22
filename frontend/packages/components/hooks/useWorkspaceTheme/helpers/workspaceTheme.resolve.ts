import {
  WORKSPACE_THEME_BY_PREFERENCE,
  WORKSPACE_THEME_BY_SYSTEM,
} from '../constants/workspaceTheme.resolution';

import type { ResolveWorkspaceTheme } from './types/workspaceTheme.resolve';

/** On phones follow the OS appearance without overwriting desktop preferences. */
export const resolveWorkspaceTheme: ResolveWorkspaceTheme = ({
  preference,
  systemDark,
  isMobile,
}) => {
  const preferredTheme = WORKSPACE_THEME_BY_PREFERENCE[preference];

  if (!isMobile && preferredTheme !== null) {
    return preferredTheme;
  }

  return WORKSPACE_THEME_BY_SYSTEM[String(systemDark)]!;
};
