import { ArrowLeft, CheckCircle2, Clock3, Info, LoaderCircle, Mail } from 'lucide-react';
import type { ChangeEventHandler } from 'react';
import { Link } from 'react-router-dom';

import { AuthMobileStatusArt } from '../../../AuthMobileArtwork';
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
      <form id={formId} className="rdr-login-v2__form rdr-auth-flow__form" onSubmit={onSubmit} noValidate>
        <div className="rdr-auth-flow__fields">
          <div className={`rdr-auth-flow__notice ${view.noticeSuccess ? 'rdr-auth-flow__notice--success' : 'rdr-auth-flow__notice--info'}`}
            role="status" aria-live="polite">
            <span className="rdr-auth-flow__notice-icon" aria-hidden="true"><NoticeIcon size={20}/></span>
            <span className="rdr-auth-flow__notice-copy">
              <strong>{view.noticeTitle}</strong>
              <small>{view.noticeMessage}</small>
            </span>
          </div>
          <div className="rdr-login-v2__field rdr-auth-flow__email-field">
            <label htmlFor={id}>Email address</label>
            <div className="rdr-login-v2__input-shell">
              <Mail size={18} aria-hidden="true"/>
              <input id={id} type="email" inputMode="email" autoComplete="email"
                placeholder="Email" value={state.email}
                aria-invalid={Boolean(state.emailError)}
                aria-describedby={state.emailError ? `${id}-error` : undefined}
                onChange={handleEmailChange} required/>
            </div>
            <small className={`rdr-login-v2__field-error rdr-login-v2__field-error--reserved${state.emailError ? '' : ' is-empty'}`}
              id={`${id}-error`} aria-live="polite" aria-atomic="true">
              {state.emailError || '\u00a0'}
            </small>
          </div>
        </div>
        <div className="rdr-auth-flow__actions">
          <button type="submit" className="rdr-login-v2__submit" disabled={state.busy || view.cooldown}
            aria-label={view.buttonLabel}>
            <span>{view.buttonLabel}</span>
            {state.busy && <LoaderCircle className="rdr-login-v2__spinner" size={19}/>}
            {!state.busy && view.cooldown && <Clock3 size={19}/>}
          </button>
          <div className="rdr-auth-flow__secondary-slot">
            <Link className="rdr-auth-flow__back" to="/login"><ArrowLeft size={16}/> Back to sign in</Link>
          </div>
          <div className="rdr-auth-flow__error-slot" aria-live="polite">
            {state.error && <p className="rdr-login-v2__error" role="alert">{state.error}</p>}
          </div>
        </div>
      </form>
      {view.noticeSuccess && (
        <div className="rdr-auth-flow__mobile-success" role="status" aria-live="polite">
          <AuthMobileStatusArt kind="email"/>
          <h3>Check your email</h3>
          <p>If this email address is eligible, a password reset link will be sent to <strong>{state.lastRequestedEmail}</strong>.</p>
          <Link className="rdr-login-v2__submit" to="/login">Back to sign in</Link>
          <button type="submit" form={formId} className="rdr-auth-flow__resend" disabled={state.busy || view.cooldown}>
            {view.resendLabel}
          </button>
          <button type="button" className="rdr-auth-flow__resend" onClick={onUseAnotherEmail}>Use another email</button>
        </div>
      )}
    </div>
  );
};
