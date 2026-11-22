import { createElement, forwardRef } from 'react';
import type { BoxProps } from './types/box.props';

/** Unstyled semantic layout primitive: preserves the native tag, attributes and ref. */
export const Box = forwardRef<HTMLElement, BoxProps>(({ as = 'div', dataTestId, children, ...props }, ref) =>
  createElement(as, { ...props, ref, 'data-testid': dataTestId }, children),
);
Box.displayName = 'Box';
