import { ThemePreference } from '../../types/useWorkspaceTheme';
import { readThemePreference } from '../workspaceTheme.readPreference';

describe('readThemePreference', () => {
  beforeEach(() => {
    window.localStorage.clear();
  });

  it('defaults to system for missing or unknown values', () => {
    expect(readThemePreference('app-one.theme')).toBe(ThemePreference.SYSTEM);
    window.localStorage.setItem('app-one.theme', 'invalid');
    expect(readThemePreference('app-one.theme')).toBe(ThemePreference.SYSTEM);
  });

  it('reads a supported preference', () => {
    window.localStorage.setItem('app-one.theme', ThemePreference.DARK);
    expect(readThemePreference('app-one.theme')).toBe(ThemePreference.DARK);
  });

  it('reads only the requested application namespace', () => {
    window.localStorage.setItem('app-one.theme', ThemePreference.DARK);
    expect(readThemePreference('app-two.theme')).toBe(ThemePreference.SYSTEM);
    window.localStorage.setItem('app-two.theme', ThemePreference.LIGHT);
    expect(readThemePreference('app-two.theme')).toBe(ThemePreference.LIGHT);
    expect(readThemePreference('app-one.theme')).toBe(ThemePreference.DARK);
  });
});
