import { apiRequest } from '@Business/api';
import { revokeSession } from '../revokeSession';

jest.mock('@Business/api', () => ({ apiRequest: jest.fn() }));

describe('revokeSession', () => {
  it('escapes an untrusted session id and uses the protected DELETE endpoint', async () => {
    const request = jest.mocked(apiRequest);
    request.mockResolvedValue(undefined as never);
    await revokeSession('a/b?c');
    expect(request).toHaveBeenCalledWith('/identity/protected/auth/sessions/a%2Fb%3Fc', {
      method: 'DELETE',
    });
  });
});
