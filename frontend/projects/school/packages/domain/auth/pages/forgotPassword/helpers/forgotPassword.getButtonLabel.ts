import type { ForgotPasswordState } from '../types';

export const getForgotPasswordButtonLabel = (
  state: ForgotPasswordState,
  noticeSuccess: boolean,
): string => {
  if (state.busy) return __('Sending link…');
  if (state.secondsLeft > 0) {
    return __('Resend in {{seconds}}s', { seconds: state.secondsLeft });
  }
  if (noticeSuccess) return __('Resend reset link');
  return __('Send reset link');
};
