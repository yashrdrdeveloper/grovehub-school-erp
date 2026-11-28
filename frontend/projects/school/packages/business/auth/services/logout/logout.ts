import { apiRequest } from '@Business/api';

import { AuthOperation } from '../constants/auth.operations';
import { getAuthEndpoint } from '../helpers/auth.getEndpoint';

import type { Logout } from './types/logout';
import type { StatusResult } from './types/logout.response';

/** Single-purpose, typed Identity operation. API owns its session cookies. */
export const logout: Logout = () => {
  const endpoint = getAuthEndpoint(AuthOperation.LOGOUT);

  return apiRequest<StatusResult>(endpoint.path, {
    method: endpoint.method,
  });
};
