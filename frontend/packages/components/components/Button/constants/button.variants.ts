import { cva } from 'class-variance-authority';

import { BUTTON_VARIANT } from './button.variant';
import { BUTTON_SIZE } from './button.size';
import type { ButtonVariant } from './types/button.variant';
import type { ButtonSize } from './types/button.size';

const variants = {
  [BUTTON_VARIANT.PRIMARY]: 'button--primary bg-primary text-primary-foreground hover:bg-primary/90',
  [BUTTON_VARIANT.SECONDARY]: 'button--secondary bg-secondary text-secondary-foreground hover:bg-secondary/85',
  [BUTTON_VARIANT.OUTLINE]: 'button--outline border border-input bg-background text-foreground hover:bg-muted',
  [BUTTON_VARIANT.GHOST]: 'button--ghost hover:bg-muted hover:text-foreground',
  [BUTTON_VARIANT.DANGER]: 'button--danger bg-destructive text-white hover:bg-destructive/90',
} satisfies Record<ButtonVariant, string>;

const sizes = {
  [BUTTON_SIZE.SMALL]: 'gh-button--sm',
  [BUTTON_SIZE.MEDIUM]: 'gh-button--md',
  [BUTTON_SIZE.LARGE]: 'gh-button--lg',
} satisfies Record<ButtonSize, string>;

export const buttonVariants = cva(
  'button gh-button inline-flex items-center justify-center gap-2 rounded-md font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring/50 disabled:pointer-events-none disabled:opacity-50',
  {
    variants: {
      variant: variants,
      size: sizes,
    },
    defaultVariants: {
      variant: BUTTON_VARIANT.PRIMARY,
      size: BUTTON_SIZE.MEDIUM,
    },
  },
);
