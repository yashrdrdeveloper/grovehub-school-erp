import type { ResolveWorkspaceTheme } from './types/workspaceTheme.resolve';

/** On phones follow OS appearance without overwriting the desktop preference. */
export const resolveWorkspaceTheme: ResolveWorkspaceTheme = ({ preference, systemDark, isMobile }) => {
  if (!isMobile && preference !== 'system') return preference;
  if (systemDark) return 'dark';
  return 'light';
};
