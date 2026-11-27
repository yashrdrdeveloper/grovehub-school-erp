import { CheckCircle2, Info } from '@Components/icons';
import { FORGOT_PASSWORD_LEGACY_CLASSES } from '../constants/forgotPassword.legacyClasses';

export const getForgotPasswordNoticeAppearance = (noticeSuccess: boolean) => {
  if (noticeSuccess) {
    return {
      NoticeIcon: CheckCircle2,
      className: FORGOT_PASSWORD_LEGACY_CLASSES.NOTICE_SUCCESS,
    };
  }
  return {
    NoticeIcon: Info,
    className: FORGOT_PASSWORD_LEGACY_CLASSES.NOTICE_INFO,
  };
};
