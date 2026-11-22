import { cn } from '@Base/utilities/cn';
import type { FormComponent } from './types/form.props';

export const Form: FormComponent = ({
  layout = 'stack', spacing = 'md', dataTestId, className, children, ...props
}) => (
  <form {...props} data-testid={dataTestId}
    className={cn('gh-form', `gh-form--${layout}`, `gh-form--gap-${spacing}`, className)}>
    {children}
  </form>
);
