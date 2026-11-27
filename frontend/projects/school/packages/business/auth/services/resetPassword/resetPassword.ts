import { apiRequest } from '@Business/api';

import { AuthOperation } from '../constants/auth.operations';
import { getAuthEndpoint } from '../helpers/auth.getEndpoint';

import type { ResetPassword } from './types/resetPassword';
import type { StatusResult } from './types/resetPassword.response';

/** Single-purpose, typed Identity operation. API owns its session cookies. */
export const resetPassword: ResetPassword = (newPassword) => {
  const endpoint = getAuthEndpoint(AuthOperation.RESET_PASSWORD);

  return apiRequest<StatusResult>(endpoint.path, {
    method: endpoint.method,
    body: { newPassword },
  });
};
