import { getPinnedAppClassName } from '../pinnedNavigation.getAppClassName';

describe('getPinnedAppClassName', () => {
  it('includes the selected modifier only for the active protected route', () => {
    expect(getPinnedAppClassName({ isActive: true })).toContain('active');
    expect(getPinnedAppClassName({ isActive: false })).toBe('item');
  });
});
