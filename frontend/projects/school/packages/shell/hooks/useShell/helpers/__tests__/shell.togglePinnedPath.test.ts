import { togglePinnedPath } from '../shell.togglePinnedPath';

describe('togglePinnedPath', () => {
  it('adds a missing app without mutating the original', () => {
    const pinned = Object.freeze(['/dashboard']);
    expect(togglePinnedPath(pinned, '/reports')).toEqual(['/dashboard', '/reports']);
    expect(pinned).toEqual(['/dashboard']);
  });
  it('removes an existing app without changing others', () => {
    expect(togglePinnedPath(['/dashboard', '/reports'], '/dashboard')).toEqual(['/reports']);
  });
});
