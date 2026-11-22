import { cva } from 'class-variance-authority';

export const buttonVariants = cva('button gh-button inline-flex items-center justify-center gap-2 rounded-md font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring/50 disabled:pointer-events-none disabled:opacity-50', {
  variants: {
    variant: {
      primary: 'button--primary bg-primary text-primary-foreground hover:bg-primary/90',
      secondary: 'button--secondary bg-secondary text-secondary-foreground hover:bg-secondary/85',
      outline: 'button--outline border border-input bg-background text-foreground hover:bg-muted',
      ghost: 'button--ghost hover:bg-muted hover:text-foreground',
      danger: 'button--danger bg-destructive text-white hover:bg-destructive/90',
    },
    size: { sm: 'gh-button--sm', md: 'gh-button--md', lg: 'gh-button--lg' },
  },
  defaultVariants: { variant: 'primary', size: 'md' },
});
