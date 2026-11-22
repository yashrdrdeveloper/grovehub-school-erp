export { AUTH_ROUTE } from './constants/auth.routes';
export { AUTH_ROUTE_IDS } from '../constants/auth.routeIds';
export { AUTH_ROUTE_CONFIG } from '../constants/auth.routeConfig';
export {
  getAuthRoute,
  getLoginRoute,
  getForgotPasswordRoute,
  getResetPasswordRoute,
  getSetPasswordRoute,
} from '../helpers/auth.routes';

export type { AuthRouteId } from '../constants/types/auth.routeIds';
