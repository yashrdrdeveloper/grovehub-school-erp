import { AuthOperation } from '../../constants/auth.operations';
import { AUTH_ENDPOINT } from '../../constants/auth.endpoints';
import { AUTHENTICATION_OPERATIONS } from '../auth.operations';

describe('AUTHENTICATION_OPERATIONS', () => {
  it('has a typed handler for every configured operation', () => {
    for (const operation of Object.values(AuthOperation)) {
      expect(typeof AUTHENTICATION_OPERATIONS[operation]).toBe('function');
      expect(AUTH_ENDPOINT[operation].path).toMatch(/^\/identity\//);
    }
  });

  it('does not expose an unauthenticated logout endpoint', () => {
    expect(AUTH_ENDPOINT[AuthOperation.LOGOUT].path).toBe(
      '/identity/protected/auth/logout',
    );
  });
});
