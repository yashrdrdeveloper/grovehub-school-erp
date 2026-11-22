import { cn } from '@Base/utilities/cn';
import type { FormFieldComponent } from './types/formField.props';

export const FormField: FormFieldComponent = ({
  htmlFor, label, hint, error, required = false, children, className, dataTestId, ...props
}) => (
  <div className={cn('field gh-form-field grid gap-1.5', className)} data-testid={(dataTestId ?? htmlFor) + '-field'} {...props}>
    <label className="field__label text-sm font-medium text-foreground" data-testid={(dataTestId ?? htmlFor) + '-label'} htmlFor={htmlFor}>
      {label}{required && <span className="gh-field-required" aria-hidden="true"> *</span>}
    </label>
    {children}
    {(error || hint) && <span id={`${htmlFor}-description`}
      data-testid={(dataTestId ?? htmlFor) + '-message'}
      className={cn('field__hint text-xs', error && 'field__hint--error')}
      role={error ? 'alert' : undefined}>{error || hint}</span>}
  </div>
);
