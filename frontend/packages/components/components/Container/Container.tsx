import { createElement } from 'react';
import { cn } from '@Base/utilities';
import type { ContainerComponent } from './types/container.props';

export const Container: ContainerComponent = ({
  as = 'div', maxWidth = 'lg', padded = true, className, children, ...props
}) => createElement(as, {
  ...props,
  className: cn('gh-container', `gh-container--${maxWidth}`, padded && 'gh-container--padded', className),
}, children);
