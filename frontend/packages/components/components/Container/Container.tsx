import { CONTAINER_TAG } from './constants/container.tag';
import { CONTAINER_WIDTH } from './constants/container.width';
import { createElement } from 'react';
import { cn } from '@Base/utilities/cn';
import type { ContainerComponent } from './types/container.props';

export const Container: ContainerComponent = ({
  as = CONTAINER_TAG.DIV, maxWidth = CONTAINER_WIDTH.LARGE, padded = true, className, children, dataTestId, ...props
}) => createElement(as, {
  ...props,
  'data-testid': dataTestId,
  className: cn('gh-container', `gh-container--${maxWidth}`, padded && 'gh-container--padded', className),
}, children);
