import { AlertCircle, Eye, EyeOff, LoaderCircle, LockKeyhole, Mail } from 'lucide-react';
import type { ChangeEventHandler } from 'react';
import { Button } from '@Components/components/Button';
import { Form } from '@Components/components/Form';

import { CredentialField } from '../../../../components/credentialField';
import { LOGIN_COPY } from '../../constants/login.copy';
import { LOGIN_FIELD_CONFIG } from '../../constants/login.fieldConfig';
import { LoginField } from '../../constants/login.fields';
import { LOGIN_TEST_IDS } from '../../constants/login.testIds';
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
    <Form unstyled dataTestId={LOGIN_TEST_IDS.FORM} className="rdr-login-v2__form rdr-login-v2__form--signin" onSubmit={onSubmit} noValidate>
      <CredentialField id={fieldId + '-' + LoginField.EMAIL} label={LOGIN_FIELD_CONFIG[LoginField.EMAIL].label} Icon={Mail}
        type="email" inputMode="email" autoComplete={LOGIN_FIELD_CONFIG[LoginField.EMAIL].autoComplete}
        placeholder={LOGIN_FIELD_CONFIG[LoginField.EMAIL].placeholder} value={state.email} required
        error={state.fieldErrors[LoginField.EMAIL]} onChange={handleEmailChange}/>

      {betweenFields}

      <CredentialField id={fieldId + '-' + LoginField.PASSWORD} label={LOGIN_FIELD_CONFIG[LoginField.PASSWORD].label}
        Icon={LockKeyhole} type={passwordType} groupLabel labelAction={passwordLabelAction}
        autoComplete={LOGIN_FIELD_CONFIG[LoginField.PASSWORD].autoComplete} placeholder={LOGIN_FIELD_CONFIG[LoginField.PASSWORD].placeholder}
        value={state.password} required error={state.fieldErrors[LoginField.PASSWORD]}
        onChange={handlePasswordChange}
        trailingAction={
          <button type="button" className="rdr-login-v2__visibility"
            onClick={onTogglePassword} aria-label={visibilityLabel}>
            <VisibilityIcon size={19} aria-hidden="true"/>
          </button>
        }/>

      <Button type="submit" variant="primary" size="lg" dataTestId={LOGIN_TEST_IDS.SUBMIT} className="rdr-login-v2__submit" disabled={state.isSubmitting}>
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
