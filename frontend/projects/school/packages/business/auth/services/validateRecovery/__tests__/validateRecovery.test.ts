import { apiRequest } from '@Business/api';
import { validateRecovery } from '../validateRecovery';

jest.mock('@Business/api', () => ({ apiRequest: jest.fn() }));

describe('validateRecovery', () => {
  it('supplies proof only as a recovery authorization header', async () => {
    const request = jest.mocked(apiRequest);
    request.mockResolvedValue(undefined as never);
    await validateRecovery('secret-proof');
    expect(request).toHaveBeenCalledWith('/identity/public/auth/recovery/validate', {
      method: 'GET',
      headers: { Authorization: 'Recovery secret-proof' },
    });
  });
});
