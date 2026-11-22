import { AUTH_ROUTE_IDS } from '../../constants/auth.routeIds';
import { getAuthRoute } from '../../helpers/auth.routes';

/** Compatibility path dictionary; URLs are derived exclusively from route metadata. */
export const AUTH_ROUTE = Object.freeze({
  LOGIN: getAuthRoute(AUTH_ROUTE_IDS.LOGIN).path,
  LOGIN_LEGACY: getAuthRoute(AUTH_ROUTE_IDS.LOGIN_LEGACY).path,
  FORGOT_PASSWORD: getAuthRoute(AUTH_ROUTE_IDS.FORGOT_PASSWORD).path,
  RESET_PASSWORD: getAuthRoute(AUTH_ROUTE_IDS.RESET_PASSWORD).path,
  SET_PASSWORD: getAuthRoute(AUTH_ROUTE_IDS.SET_PASSWORD).path,
  SETUP_LEGACY: getAuthRoute(AUTH_ROUTE_IDS.SETUP_LEGACY).path,
});
