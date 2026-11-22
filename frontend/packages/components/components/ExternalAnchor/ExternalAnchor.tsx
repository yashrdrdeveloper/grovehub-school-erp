import { forwardRef } from 'react';
import type { ExternalAnchorProps } from './types/externalAnchor.props';

/** External/native href links remain accessible anchors. */
export const ExternalAnchor = forwardRef<HTMLAnchorElement, ExternalAnchorProps>(({ dataTestId, ...props }, ref) =>
  <a {...props} ref={ref} data-testid={dataTestId}/>,
);
ExternalAnchor.displayName = 'ExternalAnchor';
