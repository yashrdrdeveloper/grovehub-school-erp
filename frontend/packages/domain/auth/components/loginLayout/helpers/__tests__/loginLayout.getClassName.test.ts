import { getLoginLayoutClassName } from '../loginLayout.getClassName';

describe('getLoginLayoutClassName', () => {
  it('preserves desktop variants', () => {
    expect(getLoginLayoutClassName('login', false)).toBe('rdr-login-v2 rdr-login-v2--login');
  });
  it('preserves the mobile-minimal modifier', () => {
    expect(getLoginLayoutClassName('recovery', true)).toContain('rdr-login-v2--mobile-minimal');
  });
});
