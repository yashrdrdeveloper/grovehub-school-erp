import { _trim } from '@Base/utilities/trim';
import type { GetForgotPasswordViewModel } from './types/forgotPassword.getViewModel';

/** All conditional copy is prepared before JSX rendering. */
export const getForgotPasswordViewModel: GetForgotPasswordViewModel = (state) => {
  const cooldown = state.secondsLeft > 0;
  const hasRequested = Boolean(state.lastRequestedEmail);
  const emailChanged = hasRequested && _trim(state.email) !== state.lastRequestedEmail;
  const noticeSuccess = hasRequested && !emailChanged;

  let noticeTitle = 'Secure password reset';
  let noticeMessage = 'Enter your email to receive a secure password reset link.';
  if (emailChanged) {
    noticeTitle = 'Email address updated';
    noticeMessage = 'You can send a link to the corrected address once the countdown ends.';
  }
  if (noticeSuccess) {
    noticeTitle = 'Request received!';
    noticeMessage = `If eligible, a password reset link will be sent to ${state.lastRequestedEmail}.`;
  }
  let buttonLabel = 'Send reset link';
  if (noticeSuccess) buttonLabel = 'Resend reset link';
  if (cooldown) buttonLabel = `Resend in ${state.secondsLeft}s`;
  if (state.busy) buttonLabel = 'Sending link…';
  let resendLabel = 'Resend reset link';
  if (cooldown) resendLabel = `Resend in ${state.secondsLeft}s`;
  if (state.busy) resendLabel = 'Sending…';

  return { cooldown, noticeSuccess, noticeTitle, noticeMessage, buttonLabel, resendLabel };
};
