import { getLoginPasswordViewModel } from '../login.getPasswordViewModel';

describe('getLoginPasswordViewModel', () => {
  it('does not expose password characters by default', () => {
    expect(getLoginPasswordViewModel(false)).toMatchObject({ label: 'Show password', type: 'password' });
  });
  it('returns the reversible visible presentation', () => {
    expect(getLoginPasswordViewModel(true)).toMatchObject({ label: 'Hide password', type: 'text' });
  });
});
