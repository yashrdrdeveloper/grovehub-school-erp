import type { RouteEntity } from '../../../baseEntities/routes';
import type { AuthRouteId } from '../../constants/types/auth.routeIds';

export type GetAuthRoute = (id: AuthRouteId) => RouteEntity<AuthRouteId>;
