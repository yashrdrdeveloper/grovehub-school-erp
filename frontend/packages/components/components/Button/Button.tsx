import type { ButtonHTMLAttributes } from 'react';

import { cn } from '@Base/utilities';

export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'ghost' | 'danger';
  loading?: boolean;
}

export const Button = ({
  variant = 'primary',
  loading = false,
  className,
  children,
  disabled,
  ...props
}: ButtonProps) => (
  <button
    className={cn('button', `button--${variant}`, className)}
    disabled={disabled || loading}
    {...props}
  >
    {loading ? 'Please wait…' : children}
  </button>
);
