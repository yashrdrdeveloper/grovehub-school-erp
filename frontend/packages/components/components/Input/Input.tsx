import { forwardRef, useId } from 'react';
import { cn } from '@Base/utilities/cn';
import { FormField } from '../FormField';
import { InputControl } from '../InputControl';
import type { InputComponent, InputProps } from './types/input.props';

export const Input: InputComponent = forwardRef<HTMLInputElement, InputProps>(({
  label, error, hint, className, id, required, dataTestId, controlSize = 'md', 'aria-describedby': describedBy, ...props
}, ref) => {
  const generatedId = useId();
  const inputId = id ?? generatedId;
  const descriptionId = error || hint ? `${inputId}-description` : undefined;
  const controlTestId = dataTestId ?? id ?? label.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, '');
  const combinedDescription = [describedBy, descriptionId].filter(Boolean).join(' ') || undefined;

  return (
    <FormField label={label} htmlFor={inputId} error={error} hint={hint} required={required} dataTestId={controlTestId}>
      <InputControl ref={ref} id={inputId} required={required} dataTestId={controlTestId} controlSize={controlSize}
        className={cn('field__input gh-input', error && 'field__input--error', className)}
        aria-invalid={Boolean(error)} aria-describedby={combinedDescription} {...props}/>
    </FormField>
  );
});
Input.displayName = 'Input';
