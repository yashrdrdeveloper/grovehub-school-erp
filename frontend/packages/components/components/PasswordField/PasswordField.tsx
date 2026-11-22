import { Eye, EyeOff } from 'lucide-react';
import { forwardRef, useId, useState } from 'react';
import { cn } from '@Base/utilities';
import { FormField } from '../FormField';
import type { PasswordFieldProps } from './types/passwordField.props';

export const PasswordField = forwardRef<HTMLInputElement, PasswordFieldProps>(({
  label, hint, error, initiallyVisible = false, className, id, disabled, required,
  'aria-describedby': describedBy, ...props
}, ref) => {
  const [visible, setVisible] = useState(initiallyVisible);
  const generatedId = useId();
  const inputId = id ?? generatedId;
  const descriptionId = error || hint ? `${inputId}-description` : undefined;
  const combinedDescription = [describedBy, descriptionId].filter(Boolean).join(' ') || undefined;
  const VisibilityIcon = visible ? EyeOff : Eye;
  const toggleLabel = visible ? 'Hide password' : 'Show password';
  const handleVisibilityChange = () => setVisible(current => !current);

  return (
    <FormField label={label} htmlFor={inputId} hint={hint} error={error} required={required}>
      <div className="gh-password-field">
        <input {...props} ref={ref} id={inputId} type={visible ? 'text' : 'password'}
          required={required} disabled={disabled}
          className={cn('field__input gh-input gh-password-field__input',
            error && 'field__input--error', className)}
          aria-invalid={Boolean(error)} aria-describedby={combinedDescription}/>
        <button type="button" className="gh-password-field__toggle" disabled={disabled}
          aria-label={toggleLabel} aria-pressed={visible} onClick={handleVisibilityChange}>
          <VisibilityIcon size={18} aria-hidden="true"/>
        </button>
      </div>
    </FormField>
  );
});
PasswordField.displayName = 'PasswordField';
