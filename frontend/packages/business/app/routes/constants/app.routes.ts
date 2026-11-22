import { APP_ROUTE_IDS } from '../../constants/app.routeIds';
import { getAppRoute } from '../../helpers/app.routes';

export const APP_ROUTE = Object.freeze({
  HOME: getAppRoute(APP_ROUTE_IDS.HOME).path,
  DASHBOARD: getAppRoute(APP_ROUTE_IDS.DASHBOARD).path,
});
