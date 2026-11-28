import { APP_ROUTE_IDS } from '../../constants/app.routeIds';
import { getAppRoute } from '../../helpers/app.routes';

export const APP_ROUTE = Object.freeze({
  HOME: getAppRoute(APP_ROUTE_IDS.HOME).path,
  DASHBOARD: getAppRoute(APP_ROUTE_IDS.DASHBOARD).path,
  AUDIT: getAppRoute(APP_ROUTE_IDS.AUDIT).path,
  RBAC: getAppRoute(APP_ROUTE_IDS.RBAC).path,
  USER_SETUP: getAppRoute(APP_ROUTE_IDS.USER_SETUP).path,
});
