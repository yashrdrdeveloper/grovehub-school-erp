import { cn } from '@Base/utilities/cn';
import styles from './CredentialField.module.css';
import { Box, BOX_TAG } from '@Components/components/Box';
import { getFormFieldCombinedDescription } from '@Components/components/FormField';
import { InputControl } from '@Components/components/InputControl';

import { renderCredentialFieldLabel } from './renderer/credentialField.label';

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
  const errorDescriptionId = (error && errorId) || undefined;

  const description = getFormFieldCombinedDescription(describedBy, errorDescriptionId);

  const fieldLabel = renderCredentialFieldLabel({ id, label, groupLabel, labelAction });

  return (
    <Box
      as={BOX_TAG.DIV}
      className={cn(styles['rdr-login-v2__field'], outerClassName)}
      dataTestId={dataTestId && `${dataTestId}-wrapper`}
    >
      {fieldLabel}

      <Box as={BOX_TAG.DIV} className={styles['rdr-login-v2__input-shell']}>
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
          styles['rdr-login-v2__field-error--reserved'],
          !error && styles['is-empty'],
        )}
        aria-live="polite"
        aria-atomic="true"
      >
        {error || '\u00a0'}
      </Box>
    </Box>
  );
};
