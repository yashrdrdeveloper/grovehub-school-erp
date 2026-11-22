import { APP_ROUTE } from '../constants/app.routes';

describe('APP_ROUTE', () => {
  it('publishes one canonical landing and dashboard path', () => {
    expect(APP_ROUTE.HOME).toBe('/');
    expect(APP_ROUTE.DASHBOARD).toBe('/dashboard');
  });
});
