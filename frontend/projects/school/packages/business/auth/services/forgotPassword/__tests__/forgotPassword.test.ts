import { apiRequest } from '@Business/api';
import { forgotPassword } from '../forgotPassword';

jest.mock('@Business/api', () => ({ apiRequest: jest.fn() }));

describe('forgotPassword', () => {
  it('submits the email to the public Identity endpoint', async () => {
    const request = jest.mocked(apiRequest);
    request.mockResolvedValue(undefined as never);
    await forgotPassword('parent@example.test');
    expect(request).toHaveBeenCalledWith('/identity/public/auth/forgot-password', {
      method: 'POST',
      body: { email: 'parent@example.test' },
    });
  });
});
