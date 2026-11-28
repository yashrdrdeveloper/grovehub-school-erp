import { getLoginSubmitLabel } from '../login.getSubmitLabel';

describe('getLoginSubmitLabel', () => {
  it('uses the submit label for ready state', () => {
    expect(getLoginSubmitLabel(false)).toBe('Sign in');
  });
  it('announces active submission', () => {
    expect(getLoginSubmitLabel(true)).toBe('Signing in…');
  });
});
