import { forwardRef, useId } from 'react';
import { cn } from '@Base/utilities';
import { FormField } from '../FormField';
import type { InputComponent, InputProps } from './types/input.props';

export const Input: InputComponent = forwardRef<HTMLInputElement, InputProps>(({
  label, error, hint, className, id, required, 'aria-describedby': describedBy, ...props
}, ref) => {
  const generatedId = useId();
  const inputId = id ?? generatedId;
  const descriptionId = error || hint ? `${inputId}-description` : undefined;
  const combinedDescription = [describedBy, descriptionId].filter(Boolean).join(' ') || undefined;

  return (
    <FormField label={label} htmlFor={inputId} error={error} hint={hint} required={required}>
      <input ref={ref} id={inputId} required={required}
        className={cn('field__input gh-input', error && 'field__input--error', className)}
        aria-invalid={Boolean(error)} aria-describedby={combinedDescription} {...props}/>
    </FormField>
  );
});
Input.displayName = 'Input';
