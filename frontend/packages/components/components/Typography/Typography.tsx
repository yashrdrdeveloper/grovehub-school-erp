import { TYPOGRAPHY_TAG } from './constants/typography.tag';
import { TYPOGRAPHY_VARIANT } from './constants/typography.variant';
import { TYPOGRAPHY_TONE } from './constants/typography.tone';
import { createElement } from 'react';
import { cn } from '@Base/utilities/cn';
import type { TypographyComponent } from './types/typography.props';

export const Typography: TypographyComponent = ({
  as = TYPOGRAPHY_TAG.P, variant = TYPOGRAPHY_VARIANT.BODY, tone = TYPOGRAPHY_TONE.DEFAULT, size, weight, className, children, dataTestId, ...props
}) => createElement(as, {
  ...props,
  'data-testid': dataTestId,
  className: cn('gh-typography', `gh-typography--${variant}`, `gh-typography--${tone}`, size && `gh-typography--size-${size}`, weight && `gh-typography--weight-${weight}`, className),
}, children);
