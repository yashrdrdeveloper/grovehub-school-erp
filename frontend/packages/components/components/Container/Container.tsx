import { createElement } from 'react';
import { cn } from '@Base/utilities/cn';
import type { ContainerComponent } from './types/container.props';

export const Container: ContainerComponent = ({
  as = 'div', maxWidth = 'lg', padded = true, className, children, dataTestId, ...props
}) => createElement(as, {
  ...props,
  'data-testid': dataTestId,
  className: cn('gh-container', `gh-container--${maxWidth}`, padded && 'gh-container--padded', className),
}, children);
