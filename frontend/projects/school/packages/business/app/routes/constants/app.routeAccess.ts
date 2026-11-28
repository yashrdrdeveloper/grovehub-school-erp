export { RouteType as RouteAccess } from '../../../baseEntities/routes';

/** Public routes may also have an optional guest-only rendering policy. */
export enum PublicRoutePolicy {
  OPEN = 'OPEN',
  GUEST_ONLY = 'GUEST_ONLY',
}
