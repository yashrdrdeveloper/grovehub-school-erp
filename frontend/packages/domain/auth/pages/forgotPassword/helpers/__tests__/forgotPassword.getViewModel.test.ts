import { getForgotPasswordViewModel } from '../forgotPassword.getViewModel';
import type { ForgotPasswordState } from '../../types';

const state: ForgotPasswordState = {
  email: '', lastRequestedEmail: '', busy: false,
  secondsLeft: 0, error: '', emailError: '',
};

describe('getForgotPasswordViewModel', () => {
  it('shows initial guidance', () => {
    expect(getForgotPasswordViewModel(state).noticeTitle).toBe('Secure password reset');
  });
  it('shows a truthful generic success message and countdown', () => {
    const view = getForgotPasswordViewModel({ ...state, email: 'a@b.edu', lastRequestedEmail: 'a@b.edu', secondsLeft: 30 });
    expect(view.noticeSuccess).toBe(true);
    expect(view.buttonLabel).toBe('Resend in 30s');
  });
  it('allows a changed address after cooldown', () => {
    const view = getForgotPasswordViewModel({ ...state, email: 'c@b.edu', lastRequestedEmail: 'a@b.edu' });
    expect(view.noticeTitle).toBe('Email address updated');
    expect(view.buttonLabel).toBe('Send reset link');
  });
});
