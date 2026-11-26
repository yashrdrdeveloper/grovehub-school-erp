import { forwardRef, useId } from 'react';

import { cn } from '@Base/utilities/cn';

import {
  FormField,
  getFormFieldControlTestId,
  getFormFieldDescriptionId,
  getFormFieldCombinedDescription,
} from '../FormField';
import { InputControl, INPUT_CONTROL_SIZE } from '../InputControl';

import type { InputComponent, InputProps } from './types/input.props';

/** The shared form component owns semantics; Lodash-backed ID logic lives in FormField helpers. */
export const Input: InputComponent = forwardRef<HTMLInputElement, InputProps>(
  (
    {
      label,
      error,
      hint,
      className,
      id,
      required,
      dataTestId,
      controlSize = INPUT_CONTROL_SIZE.MEDIUM,
      'aria-describedby': describedBy,
      ...props
    },
    ref,
  ) => {
    const generatedId = useId();
    const inputId = id ?? generatedId;
    const controlTestId = getFormFieldControlTestId({ label, id, dataTestId });
    const descriptionId = getFormFieldDescriptionId(inputId, error, hint);
    const combinedDescription = getFormFieldCombinedDescription(describedBy, descriptionId);

    return (
      <FormField
        label={label}
        htmlFor={inputId}
        error={error}
        hint={hint}
        required={required}
        dataTestId={controlTestId}
      >
        <InputControl
          {...props}
          ref={ref}
          id={inputId}
          required={required}
          dataTestId={controlTestId}
          controlSize={controlSize}
          className={cn('field__input gh-input', error && 'field__input--error', className)}
          aria-invalid={Boolean(error)}
          aria-describedby={combinedDescription}
        />
      </FormField>
    );
  },
);

Input.displayName = 'Input';
