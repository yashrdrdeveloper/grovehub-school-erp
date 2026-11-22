import { ThemePreference } from '../../types/useWorkspaceTheme';
import { readThemePreference } from '../workspaceTheme.readPreference';

describe('readThemePreference', () => {
  beforeEach(() => {
    window.localStorage.clear();
  });

  it('defaults to system for missing or unknown values', () => {
    expect(readThemePreference()).toBe(ThemePreference.SYSTEM);
    window.localStorage.setItem('groovehub.theme', 'invalid');
    expect(readThemePreference()).toBe(ThemePreference.SYSTEM);
  });

  it('reads a supported preference', () => {
    window.localStorage.setItem('groovehub.theme', ThemePreference.DARK);
    expect(readThemePreference()).toBe(ThemePreference.DARK);
  });
});
