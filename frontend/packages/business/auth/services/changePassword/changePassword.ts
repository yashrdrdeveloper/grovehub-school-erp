import { apiRequest } from '@Business/api';

import { AuthOperation } from '../constants/auth.operations';
import { getAuthEndpoint } from '../helpers/auth.getEndpoint';

import type { ChangePassword } from './types/changePassword';
import type { StatusResult } from './types/changePassword.response';

/** Single-purpose, typed Identity operation. API owns its session cookies. */
export const changePassword: ChangePassword = (newPassword) => {
  const endpoint = getAuthEndpoint(AuthOperation.CHANGE_PASSWORD);

  return apiRequest<StatusResult>(endpoint.path, {
    method: endpoint.method,
    body: { newPassword },
  });
};
