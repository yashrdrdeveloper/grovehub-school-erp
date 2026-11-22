import { Slot } from '@radix-ui/react-slot';
import { LoaderCircle } from 'lucide-react';
import { cn } from '@Base/utilities/cn';
import { buttonVariants } from './constants/button.variants';
import type { ButtonComponent } from './types/button.props';

/** Shadcn-style CVA button. The asChild variant requires one child element. */
export const Button: ButtonComponent = ({
  variant, size, loading = false, loadingText = 'Please wait…', fullWidth = false,
  asChild = false, unstyled = false, className, children, disabled, dataTestId, type = 'button', ...props
}) => {
  const classes = cn(!unstyled && buttonVariants({ variant, size }), !unstyled && fullWidth && 'button--full', className);
  const testId = dataTestId ?? props.id;
  if (asChild) {
    return (
      <Slot data-testid={testId} className={cn(classes, (disabled || loading) && 'pointer-events-none')}
        aria-disabled={disabled || loading || undefined} aria-busy={loading || undefined} {...props}>
        {children}
      </Slot>
    );
  }
  let content = children;
  if (loading) content = <><LoaderCircle aria-hidden="true" className="gh-button__spinner" size={16}/><span>{loadingText}</span></>;
  return (
    <button type={type} className={classes} disabled={disabled || loading}
      data-testid={testId} aria-busy={loading || undefined} {...props}>
      {content}
    </button>
  );
};
