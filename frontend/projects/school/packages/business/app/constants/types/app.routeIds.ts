import type { APP_ROUTE_IDS } from '../app.routeIds';

export type AppRouteId = (typeof APP_ROUTE_IDS)[keyof typeof APP_ROUTE_IDS];
