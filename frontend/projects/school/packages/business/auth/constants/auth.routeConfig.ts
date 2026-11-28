import { RouteType } from '../../baseEntities/routes';

import { AUTH_ROUTE_IDS } from './auth.routeIds';

import type { RouteEntity } from '../../baseEntities/routes';
import type { AuthRouteId } from './types/auth.routeIds';

/** Public Identity browser destinations; each URL has one authoritative definition. */
export const AUTH_ROUTE_CONFIG = [
  { id: AUTH_ROUTE_IDS.LOGIN, name: 'Sign in', path: '/login', type: RouteType.PUBLIC },
  {
    id: AUTH_ROUTE_IDS.LOGIN_LEGACY,
    name: 'Legacy sign in',
    path: '/loginv2',
    type: RouteType.PUBLIC,
  },
  {
    id: AUTH_ROUTE_IDS.FORGOT_PASSWORD,
    name: 'Forgot password',
    path: '/forgot-password',
    type: RouteType.PUBLIC,
  },
  {
    id: AUTH_ROUTE_IDS.RESET_PASSWORD,
    name: 'Reset password',
    path: '/reset-password',
    type: RouteType.PUBLIC,
  },
  {
    id: AUTH_ROUTE_IDS.SET_PASSWORD,
    name: 'Set password',
    path: '/set-password',
    type: RouteType.PUBLIC,
  },
  {
    id: AUTH_ROUTE_IDS.SETUP_LEGACY,
    name: 'Legacy account setup',
    path: '/setup',
    type: RouteType.PUBLIC,
  },
] as const satisfies readonly RouteEntity<AuthRouteId>[];
