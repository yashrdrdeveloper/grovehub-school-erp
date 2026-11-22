import { AppLink } from '@Components/components/AppLink';
import { Box, BOX_TAG } from '@Components/components/Box';
import { Button } from '@Components/components/Button';
import { AUTH_ROUTE } from '@Business/auth/routes';

import { FORGOT_PASSWORD_LEGACY_CLASSES } from '../constants/forgotPassword.legacyClasses';
import { AuthStatusArt } from '../../../../../components/authStatusArt';
import { AUTH_STATUS_KIND } from '../../../../../components/authStatusArt/constants/authStatusArt.kind';

import type { RenderForgotPasswordSuccess } from './types/forgotPassword.success';

/** A pure status panel; submit and edit actions belong to the containing hook. */
export const renderForgotPasswordSuccess: RenderForgotPasswordSuccess = ({
  formId,
  state,
  view,
  onUseAnotherEmail,
}) => (
  <Box
    as={BOX_TAG.DIV}
    className={FORGOT_PASSWORD_LEGACY_CLASSES.SUCCESS}
    role="status"
    aria-live="polite"
  >
    <AuthStatusArt kind={AUTH_STATUS_KIND.EMAIL} />
    <Box as={BOX_TAG.H3}>Check your email</Box>
    <Box as={BOX_TAG.P}>
      If this email address is eligible, a password reset link will be sent to{' '}
      <Box as={BOX_TAG.STRONG}>{state.lastRequestedEmail}</Box>.
    </Box>
    <AppLink className={FORGOT_PASSWORD_LEGACY_CLASSES.SUBMIT} to={AUTH_ROUTE.LOGIN}>
      Back to sign in
    </AppLink>
    <Button
      unstyled
      type="submit"
      form={formId}
      className={FORGOT_PASSWORD_LEGACY_CLASSES.RESEND}
      disabled={state.busy || view.cooldown}
    >
      {view.resendLabel}
    </Button>
    <Button
      unstyled
      type="button"
      className={FORGOT_PASSWORD_LEGACY_CLASSES.RESEND}
      onClick={onUseAnotherEmail}
    >
      Use another email
    </Button>
  </Box>
);
