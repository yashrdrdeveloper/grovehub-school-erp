import { ThemePreference } from '../../types/useWorkspaceTheme';
import { writeThemePreference } from '../workspaceTheme.writePreference';

describe('writeThemePreference', () => {
  it('persists one of the named theme enum values', () => {
    writeThemePreference(ThemePreference.LIGHT, 'app-one.theme');
    expect(window.localStorage.getItem('app-one.theme')).toBe(ThemePreference.LIGHT);
  });

  it('writes only the requested application namespace', () => {
    window.localStorage.setItem('app-one.theme', ThemePreference.DARK);
    writeThemePreference(ThemePreference.LIGHT, 'app-two.theme');
    expect(window.localStorage.getItem('app-two.theme')).toBe(ThemePreference.LIGHT);
    expect(window.localStorage.getItem('app-one.theme')).toBe(ThemePreference.DARK);
  });
});
