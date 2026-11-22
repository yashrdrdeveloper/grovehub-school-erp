import { APP_ROUTE_CONFIG, RouteAccess } from '@Business/app/routes';
import { AUTH_ROUTE_CONFIG } from '@Business/auth/routes';

import { publicRoutes, workspaceRoutes } from '../shell.routes';

describe('shell.routes', () => {
  it('maps public routes from the business registry without duplicating URLs', () => {
    expect(publicRoutes.map((route) => route.path)).toEqual(
      AUTH_ROUTE_CONFIG.map((route) => route.path),
    );
    expect(publicRoutes.every((route) => Boolean(route.element))).toBe(true);
  });

  it('maps only protected app routes and preserves the workspace route values', () => {
    const protectedPaths = APP_ROUTE_CONFIG
      .filter((route) => route.type === RouteAccess.PROTECTED)
      .map((route) => route.path);

    expect(workspaceRoutes.map((route) => route.path)).toEqual(protectedPaths);
  });
});
