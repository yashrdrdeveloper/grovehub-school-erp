import type { ButtonHTMLAttributes, ReactElement } from 'react';
import type { VariantProps } from 'class-variance-authority';
import type { buttonVariants } from '../config/buttonVariants.config';

export interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement>, VariantProps<typeof buttonVariants> {
  readonly loading?: boolean;
  readonly loadingText?: string;
  readonly fullWidth?: boolean;
}
export type ButtonComponent = (props: ButtonProps) => ReactElement;
