import { AuthOperation } from '../../constants/auth.operations';
import { getAuthEndpoint } from '../auth.getEndpoint';

describe('getAuthEndpoint', () => {
  it('preserves separate public and protected endpoints', () => {
    expect(getAuthEndpoint(AuthOperation.LOGIN).path).toBe(
      '/identity/public/auth/login',
    );
    expect(getAuthEndpoint(AuthOperation.ME).path).toBe(
      '/identity/protected/auth/me',
    );
  });
});
