import { getRoute } from '../../helpers/routes';

import { APP_ROUTE_CONFIG } from '../constants/app.routeConfig';
import { APP_ROUTE_IDS } from '../constants/app.routeIds';

import type { RouteEntity } from '../../baseEntities/routes';
import type { AppRouteId } from '../constants/types/app.routeIds';

export const getAppRoute = (id: AppRouteId): RouteEntity<AppRouteId> =>
  getRoute(APP_ROUTE_CONFIG, id);

export const getHomeRoute = (): RouteEntity<AppRouteId> =>
  getAppRoute(APP_ROUTE_IDS.HOME);

export const getDashboardRoute = (): RouteEntity<AppRouteId> =>
  getAppRoute(APP_ROUTE_IDS.DASHBOARD);
