import type { ButtonHTMLAttributes, ReactElement } from 'react';
import type { VariantProps } from 'class-variance-authority';
import type { buttonVariants } from '../constants/button.variants';

export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement>, VariantProps<typeof buttonVariants> {
  readonly loading?: boolean;
  readonly loadingText?: string;
  readonly fullWidth?: boolean;
  readonly dataTestId?: string;
  /** Render a single child with Radix Slot (e.g. a React Router Link). */
  readonly asChild?: boolean;
  /** Preserve native classes during the shared component migration. */
  readonly unstyled?: boolean;
}
export type ButtonComponent = (props: ButtonProps) => ReactElement;
