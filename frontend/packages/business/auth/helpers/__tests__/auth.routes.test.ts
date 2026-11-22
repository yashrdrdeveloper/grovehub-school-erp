import { RouteType } from '../../baseEntities/routes';
import { AUTH_ROUTE_IDS } from '../constants/auth.routeIds';
import { getAuthRoute } from '../auth.routes';

describe('getAuthRoute', () => {
  it('resolves an ID to its typed public navigation destination', () => {
    expect(getAuthRoute(AUTH_ROUTE_IDS.LOGIN)).toMatchObject({
      name: 'Sign in',
      path: '/login',
      type: RouteType.PUBLIC,
    });
  });

  it('preserves the historical invitation entrypoint', () => {
    expect(getAuthRoute(AUTH_ROUTE_IDS.SETUP_LEGACY).path).toBe('/setup');
  });
});
