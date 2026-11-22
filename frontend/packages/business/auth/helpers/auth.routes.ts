import { getRoute } from '../../helpers/routes';

import { AUTH_ROUTE_CONFIG } from '../constants/auth.routeConfig';
import { AUTH_ROUTE_IDS } from '../constants/auth.routeIds';

import type { RouteEntity } from '../../baseEntities/routes';
import type { AuthRouteId } from '../constants/types/auth.routeIds';
import type { GetAuthRoute } from './types/auth.routes';

export const getAuthRoute: GetAuthRoute = (id) =>
  getRoute(AUTH_ROUTE_CONFIG, id);

export const getLoginRoute = (): RouteEntity<AuthRouteId> =>
  getAuthRoute(AUTH_ROUTE_IDS.LOGIN);

export const getForgotPasswordRoute = (): RouteEntity<AuthRouteId> =>
  getAuthRoute(AUTH_ROUTE_IDS.FORGOT_PASSWORD);

export const getResetPasswordRoute = (): RouteEntity<AuthRouteId> =>
  getAuthRoute(AUTH_ROUTE_IDS.RESET_PASSWORD);

export const getSetPasswordRoute = (): RouteEntity<AuthRouteId> =>
  getAuthRoute(AUTH_ROUTE_IDS.SET_PASSWORD);
