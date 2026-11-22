import type { GetForgotPasswordViewModel } from './types/getForgotPasswordViewModel';

/** All conditional copy is prepared before JSX rendering. */
export const getForgotPasswordViewModel: GetForgotPasswordViewModel = (state) => {
  const cooldown = state.secondsLeft > 0;
  const hasRequested = Boolean(state.lastRequestedEmail);
  const emailChanged = hasRequested && state.email.trim() !== state.lastRequestedEmail;
  const noticeSuccess = hasRequested && !emailChanged;

  const noticeTitle = noticeSuccess ? 'Request received!'
    : emailChanged ? 'Email address updated' : 'Secure password reset';
  const noticeMessage = noticeSuccess
    ? `If eligible, a password reset link will be sent to ${state.lastRequestedEmail}.`
    : emailChanged
      ? 'You can send a link to the corrected address once the countdown ends.'
      : 'Enter your email to receive a secure password reset link.';
  const buttonLabel = state.busy ? 'Sending link…'
    : cooldown ? `Resend in ${state.secondsLeft}s`
      : noticeSuccess ? 'Resend reset link' : 'Send reset link';
  const resendLabel = state.busy ? 'Sending…'
    : cooldown ? `Resend in ${state.secondsLeft}s` : 'Resend reset link';

  return { cooldown, noticeSuccess, noticeTitle, noticeMessage, buttonLabel, resendLabel };
};
