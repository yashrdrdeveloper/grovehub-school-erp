import { BOX_TAG } from '@Components/components/Box';
import { forwardRef, useCallback, useId, useState } from 'react';

import { cn } from '@Base/utilities/cn';
import { Eye, EyeOff } from '@Components/icons';

import { Box } from '../Box';
import { Button } from '../Button';
import {
  FormField,
  getFormFieldControlTestId,
  getFormFieldDescriptionId,
  getFormFieldCombinedDescription,
} from '../FormField';
import { InputControl, INPUT_CONTROL_SIZE } from '../InputControl';

import type { PasswordFieldProps } from './types/passwordField.props';

/** One reusable password field with keyboard-safe visibility controls. */
export const PasswordField = forwardRef<HTMLInputElement, PasswordFieldProps>(
  (
    {
      label,
      hint,
      error,
      initiallyVisible = false,
      controlSize = INPUT_CONTROL_SIZE.MEDIUM,
      className,
      id,
      disabled,
      required,
      dataTestId,
      'aria-describedby': describedBy,
      ...props
    },
    ref,
  ) => {
    const [visible, setVisible] = useState(initiallyVisible);
    const generatedId = useId();
    const inputId = id ?? generatedId;
    const controlTestId = getFormFieldControlTestId({ label, id, dataTestId });
    const descriptionId = getFormFieldDescriptionId(inputId, error, hint);
    const combinedDescription = getFormFieldCombinedDescription(describedBy, descriptionId);

    let VisibilityIcon = Eye;
    let toggleLabel = 'Show password';
    let passwordType: 'password' | 'text' = 'password';

    if (visible) {
      VisibilityIcon = EyeOff;
      toggleLabel = 'Hide password';
      passwordType = 'text';
    }

    const handleVisibilityChange = useCallback(() => {
      setVisible((current) => !current);
    }, []);

    return (
      <FormField
        label={label}
        htmlFor={inputId}
        hint={hint}
        error={error}
        required={required}
        dataTestId={controlTestId}
      >
        <Box as={BOX_TAG.DIV} className="gh-password-field" dataTestId={`${controlTestId}-shell`}>
          <InputControl
            {...props}
            ref={ref}
            id={inputId}
            type={passwordType}
            controlSize={controlSize}
            required={required}
            disabled={disabled}
            dataTestId={controlTestId}
            className={cn(
              'field__input gh-input gh-password-field__input',
              error && 'field__input--error',
              className,
            )}
            aria-invalid={Boolean(error)}
            aria-describedby={combinedDescription}
          />

          <Button
            unstyled
            type="button"
            className="gh-password-field__toggle"
            dataTestId={`${controlTestId}-visibility`}
            disabled={disabled}
            aria-label={toggleLabel}
            aria-pressed={visible}
            onClick={handleVisibilityChange}
          >
            <VisibilityIcon size={18} aria-hidden="true" />
          </Button>
        </Box>
      </FormField>
    );
  },
);

PasswordField.displayName = 'PasswordField';
