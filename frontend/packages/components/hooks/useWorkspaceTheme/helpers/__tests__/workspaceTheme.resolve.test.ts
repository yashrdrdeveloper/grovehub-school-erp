import { ResolvedTheme, ThemePreference } from '../../types/useWorkspaceTheme';
import { resolveWorkspaceTheme } from '../workspaceTheme.resolve';

describe('resolveWorkspaceTheme', () => {
  it('preserves a selected desktop theme', () => {
    expect(
      resolveWorkspaceTheme({
        preference: ThemePreference.LIGHT,
        systemDark: true,
        isMobile: false,
      }),
    ).toBe(ResolvedTheme.LIGHT);
  });

  it('follows the OS appearance on mobile without replacing the preference', () => {
    expect(
      resolveWorkspaceTheme({
        preference: ThemePreference.LIGHT,
        systemDark: true,
        isMobile: true,
      }),
    ).toBe(ResolvedTheme.DARK);

    expect(
      resolveWorkspaceTheme({
        preference: ThemePreference.DARK,
        systemDark: false,
        isMobile: true,
      }),
    ).toBe(ResolvedTheme.LIGHT);
  });

  it('resolves system preference on desktop', () => {
    expect(
      resolveWorkspaceTheme({
        preference: ThemePreference.SYSTEM,
        systemDark: false,
        isMobile: false,
      }),
    ).toBe(ResolvedTheme.LIGHT);
  });
});
