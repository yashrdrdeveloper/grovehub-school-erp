import { AlertCircle, Eye, EyeOff, KeyRound, LoaderCircle, LockKeyhole, ShieldCheck } from '@Components/icons';
import type { ChangeEventHandler } from 'react';
import { Button } from '@Components/components/Button';
import { Form } from '@Components/components/Form';

import { CredentialField } from '../../../../components/credentialField';
import type { RecoveryPasswordFormComponent } from './types/recoveryPasswordForm';

/** Password-reset and invitation form with no backend calls or flow orchestration. */
export const RecoveryPasswordForm: RecoveryPasswordFormComponent = (props) => {
  const { id, view, password, confirmation, visible, fieldErrors, error,
    onSubmit, onPasswordChange, onConfirmationChange, onToggleVisibility } = props;
  const handlePasswordChange: ChangeEventHandler<HTMLInputElement> = event => {
    onPasswordChange(event.currentTarget.value);
  };
  const handleConfirmationChange: ChangeEventHandler<HTMLInputElement> = event => {
    onConfirmationChange(event.currentTarget.value);
  };
  let VisibilityIcon = Eye;
  let passwordType: 'text' | 'password' = 'password';
  let visibilityLabel = 'Show passwords';
  let confirmationLabel = 'Show confirmation password';
  if (visible) {
    VisibilityIcon = EyeOff;
    passwordType = 'text';
    visibilityLabel = 'Hide passwords';
    confirmationLabel = 'Hide confirmation password';
  }
  const saving = view === 'saving';
  let buttonLabel = 'Update password';
  if (saving) buttonLabel = 'Securing your account…';

  return (
    <Form unstyled dataTestId="recovery-form" className="rdr-login-v2__form rdr-recovery__form" onSubmit={onSubmit} noValidate>
      <div className="rdr-recovery__fields">
        <CredentialField id={id + '-password'} label="New password" Icon={LockKeyhole}
          type={passwordType} autoComplete="new-password" placeholder="New password"
          value={password} minLength={12} disabled={saving} required
          error={fieldErrors.password} onChange={handlePasswordChange}
          trailingAction={
            <button type="button" className="rdr-login-v2__visibility" disabled={saving}
              aria-label={visibilityLabel} onClick={onToggleVisibility}>
              <VisibilityIcon size={19} aria-hidden="true"/>
            </button>
          }/>

        <CredentialField id={id + '-confirmation'} label="Confirm password" Icon={KeyRound}
          type={passwordType} autoComplete="new-password" placeholder="Confirm password"
          value={confirmation} disabled={saving} required
          error={fieldErrors.confirmation} onChange={handleConfirmationChange}
          trailingAction={
            <button type="button" className="rdr-login-v2__visibility" disabled={saving}
              aria-label={confirmationLabel}
              onClick={onToggleVisibility}>
              <VisibilityIcon size={19} aria-hidden="true"/>
            </button>
          }/>
        <p className="rdr-recovery__requirement">
          <ShieldCheck size={17}/> Use at least 12 characters to keep your account secure.
        </p>
      </div>
      <div className="rdr-recovery__actions">
        <Button type="submit" variant="primary" size="lg" dataTestId="recovery-submit" className="rdr-login-v2__submit" disabled={saving}>
          {buttonLabel}
          {saving && <LoaderCircle className="rdr-login-v2__spinner" size={18}/>}
        </Button>
        <span className="rdr-recovery__remaining" aria-hidden="true">{'\u00a0'}</span>
        <div className="rdr-recovery__error-slot" aria-live="polite">
          {error && <p role="alert" className="rdr-recovery__inline-error"><AlertCircle size={16}/>{error}</p>}
        </div>
      </div>
    </Form>
  );
};
