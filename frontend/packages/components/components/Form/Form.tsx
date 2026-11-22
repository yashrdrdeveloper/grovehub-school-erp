import { cn } from '@Base/utilities/cn';
import type { FormComponent } from './types/form.props';

export const Form: FormComponent = ({
  layout = 'stack', spacing = 'md', unstyled = false, dataTestId, className, children, ...props
}) => (
  <form {...props} data-testid={dataTestId}
    className={cn(!unstyled && 'gh-form', !unstyled && `gh-form--${layout}`, !unstyled && `gh-form--gap-${spacing}`, className)}>
    {children}
  </form>
);
