import { LoaderCircle } from 'lucide-react';
import { cn } from '@Base/utilities/cn';
import { buttonVariants } from './config/buttonVariants.config';
import type { ButtonComponent } from './types/button.props';

export const Button: ButtonComponent = ({
  variant, size, loading = false, loadingText = 'Please wait…', fullWidth = false,
  className, children, disabled, dataTestId, type = 'button', ...props
}) => (
  <button type={type} className={cn(buttonVariants({ variant, size }),
    fullWidth && 'button--full', className)}
    disabled={disabled || loading} data-testid={dataTestId ?? props.id ?? 'ui-button'} aria-busy={loading || undefined} {...props}>
    {loading ? <><LoaderCircle aria-hidden="true" className="gh-button__spinner" size={16}/><span>{loadingText}</span></> : children}
  </button>
);
