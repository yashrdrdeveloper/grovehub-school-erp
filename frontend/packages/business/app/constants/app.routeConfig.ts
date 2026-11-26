import { RouteType } from '../../baseEntities/routes';

import { APP_ROUTE_IDS } from './app.routeIds';

import type { RouteEntity } from '../../baseEntities/routes';
import type { AppRouteId } from './types/app.routeIds';

export const APP_ROUTE_CONFIG = [
  { id: APP_ROUTE_IDS.HOME, name: 'School home', path: '/', type: RouteType.PUBLIC },
  { id: APP_ROUTE_IDS.DASHBOARD, name: 'Workspace', path: '/dashboard', type: RouteType.PROTECTED },
  {
    id: APP_ROUTE_IDS.USER_SETUP,
    name: 'User Setup',
    path: '/users/setup',
    type: RouteType.PROTECTED,
  },
  { id: APP_ROUTE_IDS.AUDIT, name: 'Audit logs', path: '/audit/logs', type: RouteType.PROTECTED },
  { id: APP_ROUTE_IDS.RBAC, name: 'Roles and permissions', path: '/rbac/roles', type: RouteType.PROTECTED },
] as const satisfies readonly RouteEntity<AppRouteId>[];
