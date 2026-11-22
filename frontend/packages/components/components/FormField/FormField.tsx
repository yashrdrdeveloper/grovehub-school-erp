import { cn } from '@Base/utilities';
import type { FormFieldComponent } from './types/formField.props';

export const FormField: FormFieldComponent = ({
  htmlFor, label, hint, error, required = false, children, className, ...props
}) => (
  <div className={cn('field gh-form-field', className)} {...props}>
    <label className="field__label" htmlFor={htmlFor}>
      {label}{required && <span className="gh-field-required" aria-hidden="true"> *</span>}
    </label>
    {children}
    {(error || hint) && <span id={`${htmlFor}-description`}
      className={cn('field__hint', error && 'field__hint--error')}
      role={error ? 'alert' : undefined}>{error || hint}</span>}
  </div>
);
