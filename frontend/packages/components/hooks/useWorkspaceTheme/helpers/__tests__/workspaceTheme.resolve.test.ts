import { resolveWorkspaceTheme } from '../workspaceTheme.resolve';

describe('resolveWorkspaceTheme', () => {
  it('preserves the desktop preference', () => {
    expect(resolveWorkspaceTheme({ preference: 'light', systemDark: true, isMobile: false })).toBe('light');
  });
  it('uses system appearance on mobile without mutating preference', () => {
    expect(resolveWorkspaceTheme({ preference: 'light', systemDark: true, isMobile: true })).toBe('dark');
    expect(resolveWorkspaceTheme({ preference: 'dark', systemDark: false, isMobile: true })).toBe('light');
  });
});
