import type { ForgotPasswordState } from '../types';

/** A real reset request is acknowledged without disclosing account existence. */
export const getForgotPasswordNoticeCopy = (
  state: ForgotPasswordState,
  emailChanged: boolean,
  noticeSuccess: boolean,
) => {
  if (noticeSuccess) {
    return {
      noticeTitle: __('Request received!'),
      noticeMessage: __('If eligible, a password reset link will be sent to {{email}}.', {
        email: state.lastRequestedEmail,
      }),
    };
  }
  if (emailChanged) {
    return {
      noticeTitle: __('Email address updated'),
      noticeMessage: __('You can send a link to the corrected address once the countdown ends.'),
    };
  }
  return {
    noticeTitle: __('Secure password reset'),
    noticeMessage: __('Enter your email to receive a secure password reset link.'),
  };
};
