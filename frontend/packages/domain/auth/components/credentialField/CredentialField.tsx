import { cn } from '@Base/utilities/cn';
import { InputControl } from '@Components/components/InputControl';
import type { CredentialFieldComponent } from './types/credentialField.props';

/** Domain-owned auth appearance, shared across login, recovery and reset forms.
 * Keeps the original DOM hooks and reserved error row so mobile spacing stays stable.
 */
export const CredentialField: CredentialFieldComponent = ({
  id, label, Icon, error, outerClassName, groupLabel = false,
  labelAction, trailingAction, className, dataTestId, 'aria-describedby': describedBy,
  ...inputProps
}) => {
  const errorId = id + '-error';
  const descriptors = [describedBy];
  if (error) descriptors.push(errorId);
  const description = descriptors.filter(Boolean).join(' ') || undefined;
  let fieldLabel = <label htmlFor={id}>{label}</label>;
  if (groupLabel) fieldLabel = <div className="rdr-login-v2__password-label"><label htmlFor={id}>{label}</label>{labelAction}</div>;

  return (
    <div className={cn('rdr-login-v2__field', outerClassName)} data-testid={dataTestId && dataTestId + '-wrapper'}>
      {fieldLabel}
      <div className="rdr-login-v2__input-shell">
        <Icon size={18} aria-hidden="true"/>
        <InputControl unstyled {...inputProps} id={id} dataTestId={dataTestId ?? id} className={className}
          aria-invalid={Boolean(error)} aria-describedby={description}/>
        {trailingAction}
      </div>
      <small id={errorId}
        className={cn('rdr-login-v2__field-error rdr-login-v2__field-error--reserved', !error && 'is-empty')}
        aria-live="polite" aria-atomic="true">{error || '\u00a0'}</small>
    </div>
  );
};
