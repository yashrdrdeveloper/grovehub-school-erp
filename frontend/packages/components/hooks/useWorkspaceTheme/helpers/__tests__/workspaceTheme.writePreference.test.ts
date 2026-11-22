import { ThemePreference } from '../../types/useWorkspaceTheme';
import { writeThemePreference } from '../workspaceTheme.writePreference';

describe('writeThemePreference', () => {
  it('persists one of the named theme enum values', () => {
    writeThemePreference(ThemePreference.LIGHT);
    expect(window.localStorage.getItem('groovehub.theme')).toBe(ThemePreference.LIGHT);
  });
});
