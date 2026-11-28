import { APP_ROUTE_CONFIG, RouteAccess } from '@Business/app/routes';

import { workspaceRoutes } from '../shell.routes';

describe('workspaceRoutes', () => {
  it('maps only protected app routes and preserves workspace route URLs', () => {
    const protectedPaths = APP_ROUTE_CONFIG.filter(
      (route) => route.type === RouteAccess.PROTECTED,
    ).map((route) => route.path);

    expect(workspaceRoutes.map((route) => route.path)).toEqual(protectedPaths);
  });
});
