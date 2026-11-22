import { AlertCircle, LoaderCircle, LockKeyhole, Mail } from '@Components/icons';
import type { ChangeEventHandler } from 'react';
import { Button } from '@Components/components/Button';
import { Form } from '@Components/components/Form';

import { CredentialField } from '../../../../components/credentialField';
import { getLoginPasswordViewModel } from '../../helpers/login.getPasswordViewModel';
import { getLoginSubmitLabel } from '../../helpers/login.getSubmitLabel';
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
  const { Icon: VisibilityIcon, label: visibilityLabel, type: passwordType } = getLoginPasswordViewModel(state.isPasswordVisible);
  const submitLabel = getLoginSubmitLabel(state.isSubmitting);

  return (
    <Form unstyled dataTestId={LOGIN_TEST_IDS.FORM} className="rdr-login-v2__form rdr-login-v2__form--signin" onSubmit={onSubmit} noValidate>
      <CredentialField dataTestId={LOGIN_TEST_IDS.EMAIL} id={fieldId + '-' + LoginField.EMAIL} label={LOGIN_FIELD_CONFIG[LoginField.EMAIL].label} Icon={Mail}
        type="email" inputMode="email" autoComplete={LOGIN_FIELD_CONFIG[LoginField.EMAIL].autoComplete}
        placeholder={LOGIN_FIELD_CONFIG[LoginField.EMAIL].placeholder} value={state.email} required
        error={state.fieldErrors[LoginField.EMAIL]} onChange={handleEmailChange}/>

      {betweenFields}

      <CredentialField dataTestId={LOGIN_TEST_IDS.PASSWORD} id={fieldId + '-' + LoginField.PASSWORD} label={LOGIN_FIELD_CONFIG[LoginField.PASSWORD].label}
        Icon={LockKeyhole} type={passwordType} groupLabel labelAction={passwordLabelAction}
        autoComplete={LOGIN_FIELD_CONFIG[LoginField.PASSWORD].autoComplete} placeholder={LOGIN_FIELD_CONFIG[LoginField.PASSWORD].placeholder}
        value={state.password} required error={state.fieldErrors[LoginField.PASSWORD]}
        onChange={handlePasswordChange}
        trailingAction={
          <button type="button" className="rdr-login-v2__visibility" data-testid={LOGIN_TEST_IDS.VISIBILITY}
            onClick={onTogglePassword} aria-label={visibilityLabel}>
            <VisibilityIcon size={19} aria-hidden="true"/>
          </button>
        }/>

      <Button type="submit" variant="primary" size="lg" dataTestId={LOGIN_TEST_IDS.SUBMIT} className="rdr-login-v2__submit" disabled={state.isSubmitting}>
        <span>{submitLabel}</span>
        {state.isSubmitting && <LoaderCircle className="rdr-login-v2__spinner" size={19} aria-hidden="true"/>}
      </Button>
      <div className="rdr-login-v2__feedback-slot rdr-login-v2__feedback-slot--signin" aria-live="polite" data-testid={LOGIN_TEST_IDS.FEEDBACK}>
        {state.errorMessage && <p className="rdr-login-v2__auth-error" role="alert" data-testid={LOGIN_TEST_IDS.AUTH_ERROR}>
          <AlertCircle size={16} aria-hidden="true"/>
          <span>{state.errorMessage}</span>
        </p>}
      </div>
    </Form>
  );
};
