import { forwardRef } from 'react';
import { Link } from 'react-router-dom';
import type { AppLinkProps } from './types/appLink.props';

/** Router-aware native anchor; never replace navigation with a button. */
export const AppLink = forwardRef<HTMLAnchorElement, AppLinkProps>(({ dataTestId, ...props }, ref) =>
  <Link {...props} ref={ref} data-testid={dataTestId}/>,
);
AppLink.displayName = 'AppLink';
