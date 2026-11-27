import { _trim } from '@Base/utilities/trim';
import type { GetForgotPasswordViewModel } from './types/forgotPassword.getViewModel';
import { getForgotPasswordNoticeCopy } from './forgotPassword.getNoticeCopy';
import { getForgotPasswordButtonLabel } from './forgotPassword.getButtonLabel';
import { getForgotPasswordResendLabel } from './forgotPassword.getResendLabel';

/** All conditional labels are produced by pure, separately tested message selectors. */
export const getForgotPasswordViewModel: GetForgotPasswordViewModel = (state) => {
  const cooldown = state.secondsLeft > 0;
  const hasRequested = Boolean(state.lastRequestedEmail);
  const emailChanged = hasRequested && _trim(state.email) !== state.lastRequestedEmail;
  const noticeSuccess = hasRequested && !emailChanged;
  const { noticeTitle, noticeMessage } = getForgotPasswordNoticeCopy(
    state,
    emailChanged,
    noticeSuccess,
  );
  const buttonLabel = getForgotPasswordButtonLabel(state, noticeSuccess);
  const resendLabel = getForgotPasswordResendLabel(state);

  return { cooldown, noticeSuccess, noticeTitle, noticeMessage, buttonLabel, resendLabel };
};
