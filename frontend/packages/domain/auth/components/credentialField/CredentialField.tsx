import { cn } from '@Base/utilities/cn';
import { Box, BOX_TAG } from '@Components/components/Box';
import { FieldLabel } from '@Components/components/FieldLabel';
import { getFormFieldCombinedDescription } from '@Components/components/FormField';
import { InputControl } from '@Components/components/InputControl';

import type { CredentialFieldComponent } from './types/credentialField.props';

/** Preserve one reserved auth error row while delegating ARIA joins to shared FormField. */
export const CredentialField: CredentialFieldComponent = ({
  id,
  label,
  Icon,
  error,
  outerClassName,
  groupLabel = false,
  labelAction,
  trailingAction,
  className,
  dataTestId,
  'aria-describedby': describedBy,
  ...inputProps
}) => {
  const errorId = `${id}-error`;
  let errorDescriptionId: string | undefined;

  if (error) {
    errorDescriptionId = errorId;
  }

  const description = getFormFieldCombinedDescription(describedBy, errorDescriptionId);

  let fieldLabel = <FieldLabel htmlFor={id}>{label}</FieldLabel>;

  if (groupLabel) {
    fieldLabel = (
      <Box as={BOX_TAG.DIV} className="rdr-login-v2__password-label">
        <FieldLabel htmlFor={id}>{label}</FieldLabel>
        {labelAction}
      </Box>
    );
  }

  return (
    <Box
      as={BOX_TAG.DIV}
      className={cn('rdr-login-v2__field', outerClassName)}
      dataTestId={dataTestId && `${dataTestId}-wrapper`}
    >
      {fieldLabel}

      <Box as={BOX_TAG.DIV} className="rdr-login-v2__input-shell">
        <Icon size={18} aria-hidden="true" />
        <InputControl
          {...inputProps}
          unstyled
          id={id}
          dataTestId={dataTestId ?? id}
          className={className}
          aria-invalid={Boolean(error)}
          aria-describedby={description}
        />
        {trailingAction}
      </Box>

      <Box
        as={BOX_TAG.SMALL}
        id={errorId}
        className={cn(
          'rdr-login-v2__field-error rdr-login-v2__field-error--reserved',
          !error && 'is-empty',
        )}
        aria-live="polite"
        aria-atomic="true"
      >
        {error || '\u00a0'}
      </Box>
    </Box>
  );
};
