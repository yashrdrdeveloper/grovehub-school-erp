import { Box, BOX_TAG } from '@Components/components/Box';
import { CheckCircle2, Info } from '@Components/icons';

import { FORGOT_PASSWORD_LEGACY_CLASSES } from '../constants/forgotPassword.legacyClasses';
import { FORGOT_PASSWORD_TEST_IDS } from '../../../constants/forgotPassword.testIds';

import type { RenderForgotPasswordNotice } from './types/forgotPassword.notice';

/** Pure feedback selection for both informational and success notice states. */
export const renderForgotPasswordNotice: RenderForgotPasswordNotice = (view) => {
  let NoticeIcon = Info;
  let className: string = FORGOT_PASSWORD_LEGACY_CLASSES.NOTICE_INFO;

  if (view.noticeSuccess) {
    NoticeIcon = CheckCircle2;
    className = FORGOT_PASSWORD_LEGACY_CLASSES.NOTICE_SUCCESS;
  }

  return (
    <Box
      as={BOX_TAG.DIV}
      className={className}
      dataTestId={FORGOT_PASSWORD_TEST_IDS.NOTICE}
      role="status"
      aria-live="polite"
    >
      <Box
        as={BOX_TAG.SPAN}
        className={FORGOT_PASSWORD_LEGACY_CLASSES.NOTICE_ICON}
        aria-hidden="true"
      >
        <NoticeIcon size={20} />
      </Box>
      <Box as={BOX_TAG.SPAN} className={FORGOT_PASSWORD_LEGACY_CLASSES.NOTICE_COPY}>
        <Box as={BOX_TAG.STRONG}>{view.noticeTitle}</Box>
        <Box as={BOX_TAG.SMALL}>{view.noticeMessage}</Box>
      </Box>
    </Box>
  );
};
