import { cn } from '@Base/utilities/cn';

import { Box } from '../Box';
import { FieldLabel } from '../FieldLabel';

import type { FormFieldComponent } from './types/formField.props';

/** FormField owns native field structure; form controls only supply state and IDs. */
export const FormField: FormFieldComponent = ({
  htmlFor,
  label,
  hint,
  error,
  required = false,
  children,
  className,
  dataTestId,
  ...props
}) => {
  const testId = dataTestId ?? htmlFor;
  let messageRole: 'alert' | undefined;

  if (error) {
    messageRole = 'alert';
  }

  return (
    <Box
      as="div"
      className={cn('field gh-form-field grid gap-1.5', className)}
      dataTestId={`${testId}-field`}
      {...props}
    >
      <FieldLabel
        className="field__label text-sm font-medium text-foreground"
        dataTestId={`${testId}-label`}
        htmlFor={htmlFor}
      >
        {label}
        {required && (
          <Box as="span" className="gh-field-required" aria-hidden="true">
            {' *'}
          </Box>
        )}
      </FieldLabel>

      {children}

      {(error || hint) && (
        <Box
          as="span"
          id={`${htmlFor}-description`}
          dataTestId={`${testId}-message`}
          className={cn('field__hint text-xs', error && 'field__hint--error')}
          role={messageRole}
        >
          {error || hint}
        </Box>
      )}
    </Box>
  );
};
