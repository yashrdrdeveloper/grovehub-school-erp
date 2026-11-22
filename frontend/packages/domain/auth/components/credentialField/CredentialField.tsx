import { Box } from '@Components/components/Box';
import { FieldLabel } from '@Components/components/FieldLabel';
import { _filter } from '@Base/utilities/filter';
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
  const description = _filter(descriptors, Boolean).join(' ') || undefined;
  let fieldLabel = <FieldLabel htmlFor={id}>{label}</FieldLabel>;
  if (groupLabel) fieldLabel = <Box as="div" className="rdr-login-v2__password-label"><FieldLabel htmlFor={id}>{label}</FieldLabel>{labelAction}</Box>;

  return (
    <Box as="div" className={cn('rdr-login-v2__field', outerClassName)} dataTestId={dataTestId && dataTestId + '-wrapper'}>
      {fieldLabel}
      <Box as="div" className="rdr-login-v2__input-shell">
        <Icon size={18} aria-hidden="true"/>
        <InputControl unstyled {...inputProps} id={id} dataTestId={dataTestId ?? id} className={className}
          aria-invalid={Boolean(error)} aria-describedby={description}/>
        {trailingAction}
      </Box>
      <small id={errorId}
        className={cn('rdr-login-v2__field-error rdr-login-v2__field-error--reserved', !error && 'is-empty')}
        aria-live="polite" aria-atomic="true">{error || '\u00a0'}</small>
    </Box>
  );
};
