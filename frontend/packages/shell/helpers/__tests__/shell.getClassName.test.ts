import { getShellClassName } from '../shell.getClassName';

describe('getShellClassName', () => {
  it('applies collapsed class only when collapsed', () => {
    expect(getShellClassName(false)).toBe('shell');
    expect(getShellClassName(true)).toBe('shell shell--collapsed');
  });
});
