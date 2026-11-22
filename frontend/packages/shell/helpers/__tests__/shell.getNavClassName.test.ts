import { getShellNavClassName } from '../shell.getNavClassName';

describe('getShellNavClassName', () => {
  it('selects the active navigation class by route state', () => {
    expect(getShellNavClassName({ isActive: true })).toBe('nav-item nav-item--active');
    expect(getShellNavClassName({ isActive: false })).toBe('nav-item');
  });
});
