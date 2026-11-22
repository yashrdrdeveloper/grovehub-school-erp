import { AlertCircle, Eye, EyeOff, LoaderCircle, LockKeyhole, Mail } from 'lucide-react';
import type { ChangeEventHandler } from 'react';
import { Button } from '@Components/components/Button';
import { Form } from '@Components/components/Form';

import { CredentialField } from '../../../../components/credentialField';
import { LOGIN_COPY, LOGIN_FIELDS } from '../../config';
import type { LoginFormComponent } from './types/loginForm';

/** Stateless credential presentation shared by desktop and mobile templates. */
export const LoginForm: LoginFormComponent = ({
  fieldId, state, onSubmit, onEmailChange, onPasswordChange, onTogglePassword,
  passwordLabelAction, betweenFields,
}) => {
  const handleEmailChange: ChangeEventHandler<HTMLInputElement> = event => {
    onEmailChange(event.currentTarget.value);
  };
  const handlePasswordChange: ChangeEventHandler<HTMLInputElement> = event => {
    onPasswordChange(event.currentTarget.value);
  };
  const VisibilityIcon = state.isPasswordVisible ? EyeOff : Eye;
  const visibilityLabel = state.isPasswordVisible ? 'Hide password' : 'Show password';
  const passwordType = state.isPasswordVisible ? 'text' : 'password';
  const submitLabel = state.isSubmitting ? LOGIN_COPY.busyButton : LOGIN_COPY.button;

  return (
    <Form unstyled dataTestId="login-form" className="rdr-login-v2__form rdr-login-v2__form--signin" onSubmit={onSubmit} noValidate>
      <CredentialField id={fieldId + '-email'} label={LOGIN_FIELDS.email.label} Icon={Mail}
        type="email" inputMode="email" autoComplete={LOGIN_FIELDS.email.autoComplete}
        placeholder={LOGIN_FIELDS.email.placeholder} value={state.email} required
        error={state.fieldErrors.email} onChange={handleEmailChange}/>

      {betweenFields}

      <CredentialField id={fieldId + '-password'} label={LOGIN_FIELDS.password.label}
        Icon={LockKeyhole} type={passwordType} groupLabel labelAction={passwordLabelAction}
        autoComplete={LOGIN_FIELDS.password.autoComplete} placeholder={LOGIN_FIELDS.password.placeholder}
        value={state.password} required error={state.fieldErrors.password}
        onChange={handlePasswordChange}
        trailingAction={
          <button type="button" className="rdr-login-v2__visibility"
            onClick={onTogglePassword} aria-label={visibilityLabel}>
            <VisibilityIcon size={19} aria-hidden="true"/>
          </button>
        }/>

      <Button type="submit" variant="primary" size="lg" dataTestId="login-submit" className="rdr-login-v2__submit" disabled={state.isSubmitting}>
        <span>{submitLabel}</span>
        {state.isSubmitting && <LoaderCircle className="rdr-login-v2__spinner" size={19} aria-hidden="true"/>}
      </Button>
      <div className="rdr-login-v2__feedback-slot rdr-login-v2__feedback-slot--signin" aria-live="polite">
        {state.errorMessage && <p className="rdr-login-v2__auth-error" role="alert">
          <AlertCircle size={16} aria-hidden="true"/>
          <span>{state.errorMessage}</span>
        </p>}
      </div>
    </Form>
  );
};
