import { Box, BOX_TAG } from '@Components/components/Box';
import { CheckCircle2, Info } from '@Components/icons';

import { FORGOT_PASSWORD_TEST_IDS } from '../../../constants/forgotPassword.testIds';

import type { RenderForgotPasswordNotice } from './types/forgotPassword.notice';

/** Pure feedback selection for both informational and success notice states. */
export const renderForgotPasswordNotice: RenderForgotPasswordNotice = (view) => {
  let NoticeIcon = Info;
  let className = 'rdr-auth-flow__notice rdr-auth-flow__notice--info';

  if (view.noticeSuccess) {
    NoticeIcon = CheckCircle2;
    className = 'rdr-auth-flow__notice rdr-auth-flow__notice--success';
  }

  return (
    <Box
      as={BOX_TAG.DIV}
      className={className}
      dataTestId={FORGOT_PASSWORD_TEST_IDS.NOTICE}
      role="status"
      aria-live="polite"
    >
      <Box as={BOX_TAG.SPAN} className="rdr-auth-flow__notice-icon" aria-hidden="true">
        <NoticeIcon size={20} />
      </Box>
      <Box as={BOX_TAG.SPAN} className="rdr-auth-flow__notice-copy">
        <Box as={BOX_TAG.STRONG}>{view.noticeTitle}</Box>
        <Box as={BOX_TAG.SMALL}>{view.noticeMessage}</Box>
      </Box>
    </Box>
  );
};
