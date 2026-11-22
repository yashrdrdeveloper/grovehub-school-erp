import { ArrowLeft, CheckCircle2, Clock3, Info, LoaderCircle, Mail } from 'lucide-react';
import type { ChangeEventHandler } from 'react';
import { Button } from '@Components/components/Button';
import { Form } from '@Components/components/Form';
import { Link } from 'react-router-dom';
import { AUTH_ROUTE } from '@Business/auth/routes';

import { AuthStatusArt } from '../../../../components/authStatusArt';
import { CredentialField } from '../../../../components/credentialField';
import type { ForgotPasswordFormComponent } from './types/forgotPasswordForm';

/** Pure presentation shared by mobile and desktop layouts. */
export const ForgotPasswordForm: ForgotPasswordFormComponent = ({
  id, formId, state, view, onSubmit, onEmailChange, onUseAnotherEmail,
}) => {
  const handleEmailChange: ChangeEventHandler<HTMLInputElement> = (event) => {
    onEmailChange(event.currentTarget.value);
  };
  const NoticeIcon = view.noticeSuccess ? CheckCircle2 : Info;

  return (
    <div className={`rdr-auth-flow rdr-auth-flow--forgot${view.noticeSuccess ? ' rdr-auth-flow--sent' : ''}`}>
      <Form unstyled dataTestId="forgot-password-form" id={formId} className="rdr-login-v2__form rdr-auth-flow__form" onSubmit={onSubmit} noValidate>
        <div className="rdr-auth-flow__fields">
          <div className={`rdr-auth-flow__notice ${view.noticeSuccess ? 'rdr-auth-flow__notice--success' : 'rdr-auth-flow__notice--info'}`}
            role="status" aria-live="polite">
            <span className="rdr-auth-flow__notice-icon" aria-hidden="true"><NoticeIcon size={20}/></span>
            <span className="rdr-auth-flow__notice-copy">
              <strong>{view.noticeTitle}</strong>
              <small>{view.noticeMessage}</small>
            </span>
          </div>
          <CredentialField id={id} label="Email address" Icon={Mail}
            outerClassName="rdr-auth-flow__email-field" type="email" inputMode="email"
            autoComplete="email" placeholder="Email" value={state.email}
            error={state.emailError} onChange={handleEmailChange} required/>
        </div>
        <div className="rdr-auth-flow__actions">
          <Button type="submit" variant="primary" size="lg" dataTestId="forgot-password-submit" className="rdr-login-v2__submit" disabled={state.busy || view.cooldown}
            aria-label={view.buttonLabel}>
            <span>{view.buttonLabel}</span>
            {state.busy && <LoaderCircle className="rdr-login-v2__spinner" size={19}/>}
            {!state.busy && view.cooldown && <Clock3 size={19}/>}
          </Button>
          <div className="rdr-auth-flow__secondary-slot">
            <Link className="rdr-auth-flow__back" to={AUTH_ROUTE.LOGIN}><ArrowLeft size={16}/> Back to sign in</Link>
          </div>
          <div className="rdr-auth-flow__error-slot" aria-live="polite">
            {state.error && <p className="rdr-login-v2__error" role="alert">{state.error}</p>}
          </div>
        </div>
      </Form>
      {view.noticeSuccess && (
        <div className="rdr-auth-flow__mobile-success" role="status" aria-live="polite">
          <AuthStatusArt kind="email"/>
          <h3>Check your email</h3>
          <p>If this email address is eligible, a password reset link will be sent to <strong>{state.lastRequestedEmail}</strong>.</p>
          <Link className="rdr-login-v2__submit" to={AUTH_ROUTE.LOGIN}>Back to sign in</Link>
          <button type="submit" form={formId} className="rdr-auth-flow__resend" disabled={state.busy || view.cooldown}>
            {view.resendLabel}
          </button>
          <button type="button" className="rdr-auth-flow__resend" onClick={onUseAnotherEmail}>Use another email</button>
        </div>
      )}
    </div>
  );
};
