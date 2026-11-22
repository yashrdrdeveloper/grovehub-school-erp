import { AppLink } from '@Components/components/AppLink';
import { Box, BOX_TAG } from '@Components/components/Box';
import { ArrowLeft, CheckCircle2, Clock3, Info, Mail } from '@Components/icons';
import { FORGOT_PASSWORD_TEST_IDS } from '../../constants/forgotPassword.testIds';
import type { ChangeEventHandler } from 'react';
import { Button, BUTTON_VARIANT, BUTTON_SIZE } from '@Components/components/Button';
import { Form } from '@Components/components/Form';
import { AUTH_ROUTE } from '@Business/auth/routes';

import { AuthStatusArt } from '../../../../components/authStatusArt';
import { CredentialField } from '../../../../components/credentialField';
import type { ForgotPasswordFormComponent } from './types/forgotPasswordForm';

/** Pure presentation shared by mobile and desktop layouts. */
export const ForgotPasswordForm: ForgotPasswordFormComponent = ({
  id,
  formId,
  state,
  view,
  onSubmit,
  onEmailChange,
  onUseAnotherEmail,
}) => {
  const handleEmailChange: ChangeEventHandler<HTMLInputElement> = (event) => {
    onEmailChange(event.currentTarget.value);
  };
  let NoticeIcon = Info;
  let rootClass = 'rdr-auth-flow rdr-auth-flow--forgot';
  let noticeClass = 'rdr-auth-flow__notice rdr-auth-flow__notice--info';
  if (view.noticeSuccess) {
    NoticeIcon = CheckCircle2;
    rootClass += ' rdr-auth-flow--sent';
    noticeClass = 'rdr-auth-flow__notice rdr-auth-flow__notice--success';
  }

  return (
    <Box as={BOX_TAG.DIV} className={rootClass} dataTestId={FORGOT_PASSWORD_TEST_IDS.ROOT}>
      <Form
        unstyled
        dataTestId={FORGOT_PASSWORD_TEST_IDS.FORM}
        id={formId}
        className="rdr-login-v2__form rdr-auth-flow__form"
        onSubmit={onSubmit}
        noValidate
      >
        <Box as={BOX_TAG.DIV} className="rdr-auth-flow__fields">
          <Box
            as={BOX_TAG.DIV}
            className={noticeClass}
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
          <CredentialField
            id={id}
            label="Email address"
            Icon={Mail}
            outerClassName="rdr-auth-flow__email-field"
            type="email"
            inputMode="email"
            autoComplete="email"
            placeholder="Email"
            value={state.email}
            error={state.emailError}
            onChange={handleEmailChange}
            required
          />
        </Box>
        <Box as={BOX_TAG.DIV} className="rdr-auth-flow__actions">
          <Button
            type="submit"
            variant={BUTTON_VARIANT.PRIMARY}
            size={BUTTON_SIZE.LARGE}
            dataTestId={FORGOT_PASSWORD_TEST_IDS.SUBMIT}
            className="rdr-login-v2__submit"
            disabled={state.busy || view.cooldown}
            loading={state.busy}
            loadingText={view.buttonLabel}
            aria-label={view.buttonLabel}
          >
            <Box as={BOX_TAG.SPAN}>{view.buttonLabel}</Box>
            {!state.busy && view.cooldown && <Clock3 size={19} />}
          </Button>
          <Box as={BOX_TAG.DIV} className="rdr-auth-flow__secondary-slot">
            <AppLink className="rdr-auth-flow__back" to={AUTH_ROUTE.LOGIN}>
              <ArrowLeft size={16} /> Back to sign in
            </AppLink>
          </Box>
          <Box as={BOX_TAG.DIV} className="rdr-auth-flow__error-slot" aria-live="polite">
            {state.error && (
              <Box as={BOX_TAG.P} className="rdr-login-v2__error" role="alert">
                {state.error}
              </Box>
            )}
          </Box>
        </Box>
      </Form>
      {view.noticeSuccess && (
        <Box
          as={BOX_TAG.DIV}
          className="rdr-auth-flow__mobile-success"
          role="status"
          aria-live="polite"
        >
          <AuthStatusArt kind="email" />
          <Box as={BOX_TAG.H3}>Check your email</Box>
          <Box as={BOX_TAG.P}>
            If this email address is eligible, a password reset link will be sent to{' '}
            <Box as={BOX_TAG.STRONG}>{state.lastRequestedEmail}</Box>.
          </Box>
          <AppLink className="rdr-login-v2__submit" to={AUTH_ROUTE.LOGIN}>
            Back to sign in
          </AppLink>
          <Button
            unstyled
            type="submit"
            form={formId}
            className="rdr-auth-flow__resend"
            disabled={state.busy || view.cooldown}
          >
            {view.resendLabel}
          </Button>
          <Button
            unstyled
            type="button"
            className="rdr-auth-flow__resend"
            onClick={onUseAnotherEmail}
          >
            Use another email
          </Button>
        </Box>
      )}
    </Box>
  );
};
