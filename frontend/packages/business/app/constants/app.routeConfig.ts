import { RouteType } from '../../baseEntities/routes';

import { APP_ROUTE_IDS } from './app.routeIds';

import type { RouteEntity } from '../../baseEntities/routes';
import type { AppRouteId } from './types/app.routeIds';

export const APP_ROUTE_CONFIG = [
  { id: APP_ROUTE_IDS.HOME, name: 'School home', path: '/', type: RouteType.PUBLIC },
  { id: APP_ROUTE_IDS.DASHBOARD, name: 'Workspace', path: '/dashboard', type: RouteType.PROTECTED },
] as const satisfies readonly RouteEntity<AppRouteId>[];
