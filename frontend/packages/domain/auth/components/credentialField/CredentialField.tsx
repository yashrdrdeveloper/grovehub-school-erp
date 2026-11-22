import { cn } from '@Base/utilities/cn';
import type { CredentialFieldComponent } from './types/credentialField.props';

/** Domain-owned auth appearance, shared across login, recovery and reset forms.
 * Keeps the original DOM hooks and reserved error row so mobile spacing stays stable.
 */
export const CredentialField: CredentialFieldComponent = ({
  id, label, Icon, error, outerClassName, groupLabel = false,
  labelAction, trailingAction, className, 'aria-describedby': describedBy,
  ...inputProps
}) => {
  const errorId = id + '-error';
  const description = [describedBy, error ? errorId : undefined].filter(Boolean).join(' ') || undefined;

  return (
    <div className={cn('rdr-login-v2__field', outerClassName)}>
      {groupLabel
        ? <div className="rdr-login-v2__password-label">
            <label htmlFor={id}>{label}</label>
            {labelAction}
          </div>
        : <label htmlFor={id}>{label}</label>}
      <div className="rdr-login-v2__input-shell">
        <Icon size={18} aria-hidden="true"/>
        <input {...inputProps} id={id} className={className}
          aria-invalid={Boolean(error)} aria-describedby={description}/>
        {trailingAction}
      </div>
      <small id={errorId}
        className={cn('rdr-login-v2__field-error rdr-login-v2__field-error--reserved', !error && 'is-empty')}
        aria-live="polite" aria-atomic="true">{error || '\u00a0'}</small>
    </div>
  );
};
