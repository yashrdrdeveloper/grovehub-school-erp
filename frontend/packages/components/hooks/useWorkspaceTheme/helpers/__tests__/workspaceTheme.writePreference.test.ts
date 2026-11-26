import { ThemePreference } from '../../types/useWorkspaceTheme';
import { writeThemePreference } from '../workspaceTheme.writePreference';

describe('writeThemePreference', () => {
  it('persists one of the named theme enum values', () => {
    writeThemePreference(ThemePreference.LIGHT);
    expect(window.localStorage.getItem('groovehub.theme')).toBe(ThemePreference.LIGHT);
  });

  it('writes only the requested application namespace', () => {
    window.localStorage.setItem('groovehub.theme', ThemePreference.DARK);
    writeThemePreference(ThemePreference.LIGHT, 'internal.theme');
    expect(window.localStorage.getItem('internal.theme')).toBe(ThemePreference.LIGHT);
    expect(window.localStorage.getItem('groovehub.theme')).toBe(ThemePreference.DARK);
  });
});
