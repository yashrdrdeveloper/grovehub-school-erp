import { forwardRef } from 'react';
import { NavLink } from 'react-router-dom';

import type { AppNavLinkProps } from './types/appNavLink.props';

/** Router-native active navigation through the shared component boundary. */
export const AppNavLink = forwardRef<HTMLAnchorElement, AppNavLinkProps>(
  ({ dataTestId, ...props }, ref) => <NavLink {...props} ref={ref} data-testid={dataTestId} />,
);

AppNavLink.displayName = 'AppNavLink';
