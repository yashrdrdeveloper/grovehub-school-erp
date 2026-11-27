import { apiRequest } from '@Business/api';
import { resetPassword } from '../resetPassword';

jest.mock('@Business/api', () => ({ apiRequest: jest.fn() }));

describe('resetPassword', () => {
  it('submits the new password to the established reset endpoint', async () => {
    const request = jest.mocked(apiRequest);
    request.mockResolvedValue(undefined as never);
    await resetPassword('new-secret-password');
    expect(request).toHaveBeenCalledWith('/identity/public/auth/reset-password', {
      method: 'POST',
      body: { newPassword: 'new-secret-password' },
    });
  });
});
