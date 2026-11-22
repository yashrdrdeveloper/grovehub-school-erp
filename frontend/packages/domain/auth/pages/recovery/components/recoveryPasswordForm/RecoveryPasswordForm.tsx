import { AlertCircle, Eye, EyeOff, KeyRound, LoaderCircle, LockKeyhole, ShieldCheck } from 'lucide-react';
import type { ChangeEventHandler } from 'react';

import type { RecoveryPasswordFormComponent } from './types/recoveryPasswordForm';

/** Reusable reset/invitation form: no backend calls, redirects or flow orchestration. */
export const RecoveryPasswordForm: RecoveryPasswordFormComponent = (props) => {
  const { id, view, password, confirmation, visible, fieldErrors, error, invitation, attempts,
    onSubmit, onPasswordChange, onConfirmationChange, onToggleVisibility } = props;
  const handlePasswordChange: ChangeEventHandler<HTMLInputElement> = (event) => {
    onPasswordChange(event.currentTarget.value);
  };
  const handleConfirmationChange: ChangeEventHandler<HTMLInputElement> = (event) => {
    onConfirmationChange(event.currentTarget.value);
  };
  const VisibilityIcon = visible ? EyeOff : Eye;
  const passwordType = visible ? 'text' : 'password';
  const saving = view === 'saving';
  const buttonLabel = saving ? 'Securing your account…' : 'Update password';
  const remaining = invitation && attempts === 1 ? '1 activation attempt remaining' : '\u00a0';

  return (
    <form className="rdr-login-v2__form rdr-recovery__form" onSubmit={onSubmit} noValidate>
      <div className="rdr-recovery__fields">
        <div className="rdr-login-v2__field">
          <label htmlFor={`${id}-password`}>New password</label>
          <div className="rdr-login-v2__input-shell">
            <LockKeyhole size={18} aria-hidden="true"/>
            <input id={`${id}-password`} type={passwordType} autoComplete="new-password"
              placeholder="New password" value={password} minLength={12} disabled={saving}
              aria-invalid={Boolean(fieldErrors.password)}
              aria-describedby={fieldErrors.password ? `${id}-password-error` : undefined}
              onChange={handlePasswordChange} required/>
            <button type="button" className="rdr-login-v2__visibility"
              aria-label={visible ? 'Hide passwords' : 'Show passwords'}
              onClick={onToggleVisibility} disabled={saving}>
              <VisibilityIcon size={19} aria-hidden="true"/>
            </button>
          </div>
          <small id={`${id}-password-error`}
            className={`rdr-login-v2__field-error rdr-login-v2__field-error--reserved${fieldErrors.password ? '' : ' is-empty'}`}
            aria-live="polite" aria-atomic="true">{fieldErrors.password || '\u00a0'}</small>
        </div>
        <div className="rdr-login-v2__field">
          <label htmlFor={`${id}-confirmation`}>Confirm password</label>
          <div className="rdr-login-v2__input-shell">
            <KeyRound size={18} aria-hidden="true"/>
            <input id={`${id}-confirmation`} type={passwordType} autoComplete="new-password"
              placeholder="Confirm password" value={confirmation} disabled={saving}
              aria-invalid={Boolean(fieldErrors.confirmation)}
              aria-describedby={fieldErrors.confirmation ? `${id}-confirmation-error` : undefined}
              onChange={handleConfirmationChange} required/>
            <button type="button" className="rdr-login-v2__visibility"
              aria-label={visible ? 'Hide confirmation password' : 'Show confirmation password'}
              onClick={onToggleVisibility} disabled={saving}>
              <VisibilityIcon size={19} aria-hidden="true"/>
            </button>
          </div>
          <small id={`${id}-confirmation-error`}
            className={`rdr-login-v2__field-error rdr-login-v2__field-error--reserved${fieldErrors.confirmation ? '' : ' is-empty'}`}
            aria-live="polite" aria-atomic="true">{fieldErrors.confirmation || '\u00a0'}</small>
        </div>
        <p className="rdr-recovery__requirement"><ShieldCheck size={17}/> Use at least 12 characters to keep your account secure.</p>
      </div>
      <div className="rdr-recovery__actions">
        <button type="submit" className="rdr-login-v2__submit" disabled={saving}>
          {buttonLabel}
          {saving && <LoaderCircle className="rdr-login-v2__spinner" size={18}/>}
        </button>
        <span className="rdr-recovery__remaining">{remaining}</span>
        <div className="rdr-recovery__error-slot" aria-live="polite">
          {error && <p role="alert" className="rdr-recovery__inline-error"><AlertCircle size={16}/>{error}</p>}
        </div>
      </div>
    </form>
  );
};
