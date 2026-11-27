import type { ForgotPasswordState } from '../types';

export const getForgotPasswordResendLabel = (state: ForgotPasswordState): string => {
  if (state.busy) return __('Sending…');
  if (state.secondsLeft > 0) {
    return __('Resend in {{seconds}}s', { seconds: state.secondsLeft });
  }
  return __('Resend reset link');
};
