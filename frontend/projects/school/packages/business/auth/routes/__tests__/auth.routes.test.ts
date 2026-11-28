import { AUTH_ROUTE } from '../constants/auth.routes';

describe('AUTH_ROUTE', () => {
  it('retains canonical and backward-compatible auth entrypoints', () => {
    expect(AUTH_ROUTE.LOGIN).toBe('/login');
    expect(AUTH_ROUTE.LOGIN_LEGACY).toBe('/loginv2');
    expect(AUTH_ROUTE.FORGOT_PASSWORD).toBe('/forgot-password');
    expect(AUTH_ROUTE.SETUP_LEGACY).toBe('/setup');
  });
});
