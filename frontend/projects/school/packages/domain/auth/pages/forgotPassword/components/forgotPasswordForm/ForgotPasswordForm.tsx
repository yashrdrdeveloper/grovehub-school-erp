import { AppLink } from '@Components/components/AppLink';
import { Box, BOX_TAG } from '@Components/components/Box';
import { ArrowLeft, Clock3, Mail } from '@Components/icons';
import { FORGOT_PASSWORD_TEST_IDS } from '../../constants/forgotPassword.testIds';
import type { ChangeEventHandler } from 'react';
import { Button, BUTTON_VARIANT, BUTTON_SIZE } from '@Components/components/Button';
import { Form } from '@Components/components/Form';
import { AUTH_ROUTE } from '@Business/auth/routes';

import { CredentialField } from '../../../../components/credentialField';
import { renderForgotPasswordNotice } from './renderer/forgotPassword.notice';
import { renderForgotPasswordSuccess } from './renderer/forgotPassword.success';
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
  const rootClass =
    'rdr-auth-flow rdr-auth-flow--forgot' + ((view.noticeSuccess && ' rdr-auth-flow--sent') || '');

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
          {renderForgotPasswordNotice(view)}
          <CredentialField
            id={id}
            label={__('Email address')}
            Icon={Mail}
            outerClassName="rdr-auth-flow__email-field"
            type="email"
            inputMode="email"
            autoComplete="email"
            placeholder={__('Email')}
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
              <ArrowLeft size={16} />
              {__('Back to sign in')}
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
      {view.noticeSuccess &&
        renderForgotPasswordSuccess({ formId, state, view, onUseAnotherEmail })}
    </Box>
  );
};
