import { getLoginErrorMessage } from '../login.getErrorMessage';

describe('getLoginErrorMessage', () => {
  it('does not expose credential failure details', () => {
    expect(getLoginErrorMessage(new Error('Invalid credentials'))).toBe(
      'Incorrect email or password. Please try again.',
    );
    expect(getLoginErrorMessage(new Error('HTTP 401 Unauthorized'))).toBe(
      'Incorrect email or password. Please try again.',
    );
  });

  it('does not expose unknown exceptions or backend messages', () => {
    expect(getLoginErrorMessage(new Error('database 500: secret token'))).toBe(
      'Unable to sign in right now. Please try again.',
    );
    expect(getLoginErrorMessage('internal error')).toBe(
      'Unable to sign in right now. Please try again.',
    );
  });
});
