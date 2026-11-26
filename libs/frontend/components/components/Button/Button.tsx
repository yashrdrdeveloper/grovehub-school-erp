import { Slot } from '@radix-ui/react-slot';

import { cn } from '@Base/utilities/cn';

import { BUTTON_VARIANT } from './constants/button.variant';
import { BUTTON_SIZE } from './constants/button.size';
import { buttonVariants } from './constants/button.variants';
import { renderButtonContent } from './renderer/button.content';
import type { ButtonComponent } from './types/button.props';

/** Shadcn-compatible button; all variant names and sizes share one typed contract. */
export const Button: ButtonComponent = ({
  variant = BUTTON_VARIANT.PRIMARY,
  size = BUTTON_SIZE.MEDIUM,
  loading = false,
  loadingText = 'Please wait…',
  fullWidth = false,
  asChild = false,
  unstyled = false,
  className,
  children,
  disabled,
  dataTestId,
  type = 'button',
  ...props
}) => {
  const classes = cn(
    !unstyled && buttonVariants({ variant, size }),
    !unstyled && fullWidth && 'button--full',
    className,
  );
  const testId = dataTestId ?? props.id;

  if (asChild) {
    return (
      <Slot
        {...props}
        data-testid={testId}
        className={cn(classes, (disabled || loading) && 'pointer-events-none')}
        aria-disabled={disabled || loading || undefined}
        aria-busy={loading || undefined}
      >
        {children}
      </Slot>
    );
  }

  return (
    <button
      {...props}
      type={type}
      className={classes}
      disabled={disabled || loading}
      data-testid={testId}
      aria-busy={loading || undefined}
    >
      {renderButtonContent({ loading, loadingText, children })}
    </button>
  );
};
