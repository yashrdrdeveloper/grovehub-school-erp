import { Slot } from '@radix-ui/react-slot';
import { LoaderCircle } from 'lucide-react';
import { cn } from '@Base/utilities/cn';
import { buttonVariants } from './config/buttonVariants.config';
import type { ButtonComponent } from './types/button.props';

/** Shadcn-style CVA button. The asChild variant requires one child element. */
export const Button: ButtonComponent = ({
  variant, size, loading = false, loadingText = 'Please wait…', fullWidth = false,
  asChild = false, className, children, disabled, dataTestId, type = 'button', ...props
}) => {
  const classes = cn(buttonVariants({ variant, size }), fullWidth && 'button--full', className);
  const testId = dataTestId ?? props.id ?? 'ui-button';
  if (asChild) {
    return (
      <Slot data-testid={testId} className={cn(classes, (disabled || loading) && 'pointer-events-none')}
        aria-disabled={disabled || loading || undefined} aria-busy={loading || undefined} {...props}>
        {children}
      </Slot>
    );
  }
  return (
    <button type={type} className={classes} disabled={disabled || loading}
      data-testid={testId} aria-busy={loading || undefined} {...props}>
      {loading ? <><LoaderCircle aria-hidden="true" className="gh-button__spinner" size={16}/><span>{loadingText}</span></> : children}
    </button>
  );
};
