import { getForgotPasswordResendLabel } from '../forgotPassword.getResendLabel';
import type { ForgotPasswordState } from '../../types';

const state: ForgotPasswordState = {
  email: 'a@b.edu',
  lastRequestedEmail: '',
  busy: false,
  secondsLeft: 0,
  error: '',
  emailError: '',
};

describe('getForgotPasswordResendLabel', () => {
  it('prioritizes busy, countdown and default', () => {
    expect(getForgotPasswordResendLabel(state)).toBe('Resend reset link');
    expect(getForgotPasswordResendLabel({ ...state, secondsLeft: 1 })).toBe('Resend in 1s');
    expect(getForgotPasswordResendLabel({ ...state, busy: true })).toBe('Sending…');
  });
});
