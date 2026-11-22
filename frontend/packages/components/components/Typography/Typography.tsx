import { createElement } from 'react';
import { cn } from '@Base/utilities';
import type { TypographyComponent } from './types/typography.props';

export const Typography: TypographyComponent = ({
  as = 'p', variant = 'body', tone = 'default', className, children, ...props
}) => createElement(as, {
  ...props,
  className: cn('gh-typography', `gh-typography--${variant}`, `gh-typography--${tone}`, className),
}, children);
