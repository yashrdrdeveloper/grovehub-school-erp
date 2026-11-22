import { FORM_LAYOUT } from './constants/form.layout';
import { FORM_SPACING } from './constants/form.spacing';
import { cn } from '@Base/utilities/cn';
import type { FormComponent } from './types/form.props';

export const Form: FormComponent = ({
  layout = FORM_LAYOUT.STACK, spacing = FORM_SPACING.MEDIUM, unstyled = false, dataTestId, className, children, ...props
}) => (
  <form {...props} data-testid={dataTestId}
    className={cn(!unstyled && 'gh-form', !unstyled && `gh-form--${layout}`, !unstyled && `gh-form--gap-${spacing}`, className)}>
    {children}
  </form>
);
