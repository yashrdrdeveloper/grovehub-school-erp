import { apiRequest } from './http';

describe('shared HTTP envelope', () => {
  const originalFetch = globalThis.fetch;

  afterEach(() => {
    globalThis.fetch = originalFetch;
    jest.clearAllMocks();
  });

  it('reads successful responses using credentialed requests', async () => {
    const request = jest.fn().mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => ({ success: true, data: { status: 'ACTIVE' } }),
    });
    globalThis.fetch = request as typeof fetch;

    await expect(apiRequest<{ status: string }>('/identity/protected/auth/me')).resolves.toEqual({
      status: 'ACTIVE',
    });
    expect(request).toHaveBeenCalledWith(
      '/identity/protected/auth/me',
      expect.objectContaining({ credentials: 'include', cache: 'no-store' }),
    );
  });

  it('maps the backend errorCode and errorMessage to a useful error', async () => {
    const request = jest.fn().mockResolvedValue({
      ok: false,
      status: 403,
      json: async () => ({
        success: false,
        errorCode: 'IDENTITY_FORBIDDEN',
        errorMessage: 'Access denied',
        httpErrorCode: 403,
        errorMetaData: {},
      }),
    });
    globalThis.fetch = request as typeof fetch;

    await expect(apiRequest('/identity/protected/accounts')).rejects.toEqual(
      expect.objectContaining({
        code: 'IDENTITY_FORBIDDEN',
        message: 'Access denied',
        status: 403,
      }),
    );
  });
});
