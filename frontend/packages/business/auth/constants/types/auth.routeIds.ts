import type { AUTH_ROUTE_IDS } from '../auth.routeIds';

export type AuthRouteId = (typeof AUTH_ROUTE_IDS)[keyof typeof AUTH_ROUTE_IDS];
