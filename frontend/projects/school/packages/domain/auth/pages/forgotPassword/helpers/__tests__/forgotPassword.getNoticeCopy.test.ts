import { getForgotPasswordNoticeCopy } from '../forgotPassword.getNoticeCopy';
import type { ForgotPasswordState } from '../../types';

const state: ForgotPasswordState = {
  email: 'a@b.edu',
  lastRequestedEmail: 'a@b.edu',
  busy: false,
  secondsLeft: 0,
  error: '',
  emailError: '',
};

describe('getForgotPasswordNoticeCopy', () => {
  it('keeps success truthful and supports named email interpolation', () => {
    expect(getForgotPasswordNoticeCopy(state, false, true).noticeMessage).toBe(
      'If eligible, a password reset link will be sent to a@b.edu.',
    );
    expect(getForgotPasswordNoticeCopy(state, true, false).noticeTitle).toBe(
      'Email address updated',
    );
    expect(getForgotPasswordNoticeCopy(state, false, false).noticeTitle).toBe(
      'Secure password reset',
    );
  });
});
