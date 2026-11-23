import { getForgotPasswordButtonLabel } from '../forgotPassword.getButtonLabel';
import type { ForgotPasswordState } from '../../types';

const state: ForgotPasswordState = {
  email: 'a@b.edu', lastRequestedEmail: '', busy: false,
  secondsLeft: 0, error: '', emailError: '',
};

describe('getForgotPasswordButtonLabel', () => {
  it('gives sending, countdown, resend and send states the right priority', () => {
    expect(getForgotPasswordButtonLabel(state, false)).toBe('Send reset link');
    expect(getForgotPasswordButtonLabel(state, true)).toBe('Resend reset link');
    expect(getForgotPasswordButtonLabel({ ...state, secondsLeft: 30 }, true)).toBe('Resend in 30s');
    expect(getForgotPasswordButtonLabel({ ...state, secondsLeft: 30, busy: true }, true)).toBe('Sending link…');
  });
});
