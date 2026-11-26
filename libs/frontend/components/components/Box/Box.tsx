import { createElement, forwardRef } from 'react';
import { BOX_TAG } from './constants/box.tag';
import type { BoxProps } from './types/box.props';

/** Unstyled semantic layout primitive: preserves the native tag, attributes and ref. */
export const Box = forwardRef<HTMLElement, BoxProps>(
  ({ as = BOX_TAG.DIV, dataTestId, children, ...props }, ref) => {
    const attributes: typeof props & { 'data-testid'?: string } = props;
    const selector = dataTestId ?? attributes['data-testid'];
    return createElement(as, { ...props, ref, 'data-testid': selector }, children);
  },
);
Box.displayName = 'Box';
