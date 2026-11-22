import { forwardRef, useId } from 'react';
import { cn } from '@Base/utilities/cn';
import { FormField } from '../FormField';
import type { InputComponent, InputProps } from './types/input.props';

export const Input: InputComponent = forwardRef<HTMLInputElement, InputProps>(({
  label, error, hint, className, id, required, dataTestId, 'aria-describedby': describedBy, ...props
}, ref) => {
  const generatedId = useId();
  const inputId = id ?? generatedId;
  const descriptionId = error || hint ? `${inputId}-description` : undefined;
  const combinedDescription = [describedBy, descriptionId].filter(Boolean).join(' ') || undefined;

  return (
    <FormField label={label} htmlFor={inputId} error={error} hint={hint} required={required} dataTestId={dataTestId ?? inputId}>
      <input ref={ref} id={inputId} required={required}
        data-testid={dataTestId ?? inputId}
        className={cn('field__input gh-input flex h-10 w-full min-w-0 rounded-md border border-input bg-background px-3 py-2 text-sm text-foreground shadow-xs placeholder:text-muted-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring/50 disabled:cursor-not-allowed disabled:opacity-50', error && 'field__input--error', className)}
        aria-invalid={Boolean(error)} aria-describedby={combinedDescription} {...props}/>
    </FormField>
  );
});
Input.displayName = 'Input';
