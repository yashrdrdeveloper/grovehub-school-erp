import { RouteType } from '../../baseEntities/routes';
import { getRoute } from '../routes';

describe('getRoute', () => {
  const routes = [
    { id: 'HOME', name: 'Home', path: '/', type: RouteType.PUBLIC },
  ] as const;

  it('finds a typed configured browser route', () => {
    expect(getRoute(routes, 'HOME')).toMatchObject({
      id: 'HOME',
      path: '/',
      type: RouteType.PUBLIC,
    });
  });

  it('throws for missing routes instead of silently selecting another', () => {
    expect(() => getRoute(routes, 'MISSING' as 'HOME')).toThrow('Unknown route id');
  });
});
