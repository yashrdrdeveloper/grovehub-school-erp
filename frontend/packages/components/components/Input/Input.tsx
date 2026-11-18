import { forwardRef, type InputHTMLAttributes } from 'react';

import { cn } from '@Base/utilities';

export interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label: string;
  error?: string;
  hint?: string;
}

export const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ label, error, hint, className, id, ...props }, ref) => (
    <label className="field" htmlFor={id}>
      <span className="field__label">{label}</span>
      <input
        ref={ref}
        id={id}
        className={cn('field__input', error && 'field__input--error', className)}
        aria-invalid={Boolean(error)}
        aria-describedby={error || hint ? `${id}-description` : undefined}
        {...props}
      />
      {error || hint ? (
        <span id={`${id}-description`} className={cn('field__hint', error && 'field__hint--error')}>
          {error || hint}
        </span>
      ) : null}
    </label>
  ),
);

Input.displayName = 'Input';
