import { Eye, EyeOff } from 'lucide-react';
import { forwardRef, useId, useState } from 'react';
import { cn } from '@Base/utilities/cn';
import { FormField } from '../FormField';
import { InputControl } from '../InputControl';
import type { PasswordFieldProps } from './types/passwordField.props';

export const PasswordField = forwardRef<HTMLInputElement, PasswordFieldProps>(({
  label, hint, error, initiallyVisible = false, controlSize = 'md', className, id, disabled, required,
  dataTestId, 'aria-describedby': describedBy, ...props
}, ref) => {
  const [visible, setVisible] = useState(initiallyVisible);
  const generatedId = useId();
  const inputId = id ?? generatedId;
  const controlTestId = dataTestId ?? id ?? label.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, '');
  const descriptionId = error || hint ? `${inputId}-description` : undefined;
  const combinedDescription = [describedBy, descriptionId].filter(Boolean).join(' ') || undefined;
  const VisibilityIcon = visible ? EyeOff : Eye;
  const toggleLabel = visible ? 'Hide password' : 'Show password';
  const handleVisibilityChange = () => setVisible(current => !current);

  return (
    <FormField label={label} htmlFor={inputId} hint={hint} error={error} required={required} dataTestId={controlTestId}>
      <div className="gh-password-field" data-testid={(controlTestId) + '-shell'}>
        <InputControl {...props} ref={ref} id={inputId} type={visible ? 'text' : 'password'} controlSize={controlSize}
          required={required} disabled={disabled} dataTestId={controlTestId}
          className={cn('field__input gh-input gh-password-field__input',
            error && 'field__input--error', className)}
          aria-invalid={Boolean(error)} aria-describedby={combinedDescription}/>
        <button type="button" className="gh-password-field__toggle" data-testid={(controlTestId) + '-visibility'} disabled={disabled}
          aria-label={toggleLabel} aria-pressed={visible} onClick={handleVisibilityChange}>
          <VisibilityIcon size={18} aria-hidden="true"/>
        </button>
      </div>
    </FormField>
  );
});
PasswordField.displayName = 'PasswordField';
