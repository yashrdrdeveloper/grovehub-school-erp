import { cva } from 'class-variance-authority';

export const buttonVariants = cva('button gh-button', {
  variants: {
    variant: {
      primary: 'button--primary',
      secondary: 'button--secondary',
      ghost: 'button--ghost',
      danger: 'button--danger',
    },
    size: { sm: 'gh-button--sm', md: 'gh-button--md', lg: 'gh-button--lg' },
  },
  defaultVariants: { variant: 'primary', size: 'md' },
});
