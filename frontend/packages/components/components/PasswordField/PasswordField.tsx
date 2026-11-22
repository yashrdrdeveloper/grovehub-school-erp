import { INPUT_CONTROL_SIZE } from '../InputControl/constants/inputControl.size';
import { Eye, EyeOff } from 'lucide-react';
import { forwardRef, useId, useState } from 'react';
import { cn } from '@Base/utilities/cn';
import { FormField } from '../FormField';
import { InputControl } from '../InputControl';
import type { PasswordFieldProps } from './types/passwordField.props';

export const PasswordField = forwardRef<HTMLInputElement, PasswordFieldProps>(({
  label, hint, error, initiallyVisible = false, controlSize = INPUT_CONTROL_SIZE.MEDIUM, className, id, disabled, required,
  dataTestId, 'aria-describedby': describedBy, ...props
}, ref) => {
  const [visible, setVisible] = useState(initiallyVisible);
  const generatedId = useId();
  const inputId = id ?? generatedId;
  const controlTestId = dataTestId ?? id ?? label.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, '');
  let descriptionId: string | undefined;
  if (error || hint) descriptionId = `${inputId}-description`;
  const combinedDescription = [describedBy, descriptionId].filter(Boolean).join(' ') || undefined;
  let VisibilityIcon = Eye;
  let toggleLabel = 'Show password';
  let passwordType: 'password' | 'text' = 'password';
  if (visible) {
    VisibilityIcon = EyeOff;
    toggleLabel = 'Hide password';
    passwordType = 'text';
  }
  const handleVisibilityChange = () => setVisible(current => !current);

  return (
    <FormField label={label} htmlFor={inputId} hint={hint} error={error} required={required} dataTestId={controlTestId}>
      <div className="gh-password-field" data-testid={(controlTestId) + '-shell'}>
        <InputControl {...props} ref={ref} id={inputId} type={passwordType} controlSize={controlSize}
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
