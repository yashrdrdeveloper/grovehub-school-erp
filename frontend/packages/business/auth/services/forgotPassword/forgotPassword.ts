import { apiRequest } from '@Business/api';

import { AuthOperation } from '../constants/auth.operations';
import { getAuthEndpoint } from '../helpers/auth.getEndpoint';

import type { ForgotPassword } from './types/forgotPassword';
import type { ForgotPasswordResponse } from './types/forgotPassword.response';

/** Single-purpose, typed Identity operation. API owns its session cookies. */
export const forgotPassword: ForgotPassword = (email) => {
  const endpoint = getAuthEndpoint(AuthOperation.FORGOT_PASSWORD);

  return apiRequest<ForgotPasswordResponse>(endpoint.path, {
    method: endpoint.method,
    body: { email },
  });
};
