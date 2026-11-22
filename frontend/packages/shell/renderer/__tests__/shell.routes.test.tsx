import { AUTH_ROUTE_CONFIG } from '@Business/auth/routes';

import { publicRoutes } from '../shell.routes';

describe('publicRoutes', () => {
  it('maps public routes from the business registry without duplicating URLs', () => {
    expect(publicRoutes.map((route) => route.path)).toEqual(
      AUTH_ROUTE_CONFIG.map((route) => route.path),
    );
    expect(publicRoutes.every((route) => Boolean(route.element))).toBe(true);
  });

});

