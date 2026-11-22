import { AlertCircle, Eye, EyeOff, LoaderCircle, LockKeyhole, Mail } from 'lucide-react';
import type { ChangeEventHandler } from 'react';

import { LOGIN_COPY, LOGIN_FIELDS } from '../../config';
import type { LoginFormComponent } from './types/loginForm';

/** Stateless credential fields used by both responsive templates. */
export const LoginForm: LoginFormComponent = ({
  fieldId, state, onSubmit, onEmailChange, onPasswordChange, onTogglePassword,
  passwordLabelAction, betweenFields,
}) => {
  const handleEmailChange: ChangeEventHandler<HTMLInputElement> = (event) => {
    onEmailChange(event.currentTarget.value);
  };
  const handlePasswordChange: ChangeEventHandler<HTMLInputElement> = (event) => {
    onPasswordChange(event.currentTarget.value);
  };
  const VisibilityIcon = state.isPasswordVisible ? EyeOff : Eye;
  const visibilityLabel = state.isPasswordVisible ? 'Hide password' : 'Show password';
  const passwordType = state.isPasswordVisible ? 'text' : 'password';
  const submitLabel = state.isSubmitting ? LOGIN_COPY.busyButton : LOGIN_COPY.button;

  return (
    <form className="rdr-login-v2__form rdr-login-v2__form--signin" onSubmit={onSubmit} noValidate>
      <div className="rdr-login-v2__field">
        <label htmlFor={`${fieldId}-email`}>{LOGIN_FIELDS.email.label}</label>
        <div className="rdr-login-v2__input-shell">
          <Mail size={18} aria-hidden="true"/>
          <input id={`${fieldId}-email`} type="email" inputMode="email"
            autoComplete={LOGIN_FIELDS.email.autoComplete} placeholder={LOGIN_FIELDS.email.placeholder}
            value={state.email} required
            aria-invalid={Boolean(state.fieldErrors.email)}
            aria-describedby={state.fieldErrors.email ? `${fieldId}-email-error` : undefined}
            onChange={handleEmailChange}/>
        </div>
        <small id={`${fieldId}-email-error`}
          className={`rdr-login-v2__field-error rdr-login-v2__field-error--reserved${state.fieldErrors.email ? '' : ' is-empty'}`}
          aria-live="polite" aria-atomic="true">
          {state.fieldErrors.email || '\u00a0'}
        </small>
      </div>

      {betweenFields}
      <div className="rdr-login-v2__field">
        <div className="rdr-login-v2__password-label">
          <label htmlFor={`${fieldId}-password`}>{LOGIN_FIELDS.password.label}</label>
          {passwordLabelAction}
        </div>
        <div className="rdr-login-v2__input-shell">
          <LockKeyhole size={18} aria-hidden="true"/>
          <input id={`${fieldId}-password`} type={passwordType}
            autoComplete={LOGIN_FIELDS.password.autoComplete}
            placeholder={LOGIN_FIELDS.password.placeholder}
            value={state.password} required
            aria-invalid={Boolean(state.fieldErrors.password)}
            aria-describedby={state.fieldErrors.password ? `${fieldId}-password-error` : undefined}
            onChange={handlePasswordChange}/>
          <button type="button" className="rdr-login-v2__visibility"
            onClick={onTogglePassword} aria-label={visibilityLabel}>
            <VisibilityIcon size={19} aria-hidden="true"/>
          </button>
        </div>
        <small id={`${fieldId}-password-error`}
          className={`rdr-login-v2__field-error rdr-login-v2__field-error--reserved${state.fieldErrors.password ? '' : ' is-empty'}`}
          aria-live="polite" aria-atomic="true">
          {state.fieldErrors.password || '\u00a0'}
        </small>
      </div>

      <button type="submit" className="rdr-login-v2__submit" disabled={state.isSubmitting}>
        <span>{submitLabel}</span>
        {state.isSubmitting && <LoaderCircle className="rdr-login-v2__spinner" size={19} aria-hidden="true"/>}
      </button>
      <div className="rdr-login-v2__feedback-slot rdr-login-v2__feedback-slot--signin" aria-live="polite">
        {state.errorMessage && <p className="rdr-login-v2__auth-error" role="alert">
          <AlertCircle size={16} aria-hidden="true"/>
          <span>{state.errorMessage}</span>
        </p>}
      </div>
    </form>
  );
};
