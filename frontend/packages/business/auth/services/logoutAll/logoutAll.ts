import { apiRequest } from '@Business/api';

import { AuthOperation } from '../../constants/auth.operations';
import { getAuthEndpoint } from '../../helpers/auth.getEndpoint';

import type { LogoutAll } from './types/logoutAll';
import type { LogoutAllResponse } from './types/logoutAll.response';

/** Single-purpose, typed Identity operation. API owns its session cookies. */
export const logoutAll: LogoutAll = () => {
  const endpoint = getAuthEndpoint(AuthOperation.LOGOUT_ALL);

  return apiRequest<LogoutAllResponse>(endpoint.path, {
    method: endpoint.method,
  });
};
