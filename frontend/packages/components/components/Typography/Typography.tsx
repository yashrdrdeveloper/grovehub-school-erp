import { createElement } from 'react';
import { cn } from '@Base/utilities/cn';
import type { TypographyComponent } from './types/typography.props';

export const Typography: TypographyComponent = ({
  as = 'p', variant = 'body', tone = 'default', size, weight, className, children, dataTestId, ...props
}) => createElement(as, {
  ...props,
  'data-testid': dataTestId,
  className: cn('gh-typography', `gh-typography--${variant}`, `gh-typography--${tone}`, size && `gh-typography--size-${size}`, weight && `gh-typography--weight-${weight}`, className),
}, children);
