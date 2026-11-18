import { ApiError } from './types';
import { apiRequest } from './http';

describe('shared HTTP envelope', () => {
  afterEach(() => jest.restoreAllMocks());

  it('reads successful responses without leaking authentication tokens', async () => {
    const spy = jest.spyOn(globalThis, 'fetch').mockResolvedValue({
      ok: true,
      status: 200,
      json: async () => ({ success: true, data: { status: 'ACTIVE' } }),
    } as Response);

    await expect(apiRequest<{ status: string }>('/identity/protected/auth/me')).resolves.toEqual({ status: 'ACTIVE' });
    expect(spy).toHaveBeenCalledWith(
      '/identity/protected/auth/me',
      expect.objectContaining({ credentials: 'include', cache: 'no-store' }),
    );
  });

  it('maps backend errorCode and errorMessage into a typed error', async () => {
    jest.spyOn(globalThis, 'fetch').mockResolvedValue({
      ok: false,
      status: 403,
      json: async () => ({
        success: false,
        errorCode: 'IDENTITY_FORBIDDEN',
        errorMessage: 'Access denied',
        httpErrorCode: 403,
        errorMetaData: {},
      }),
    } as Response);

    await expect(apiRequest('/identity/protected/accounts')).rejects.toEqual(
      expect.objectContaining<ApiError>({
        code: 'IDENTITY_FORBIDDEN',
        message: 'Access denied',
        status: 403,
      }),
    );
  });
});
