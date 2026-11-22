import { apiRequest } from '@Business/api';

import { AuthOperation } from '../../constants/auth.operations';
import { getAuthEndpoint } from '../../helpers/auth.getEndpoint';

import type { Login } from './types/login';
import type { LoginResponse } from './types/login.response';

/** Login uses the same operation metadata as every other Identity transport call. */
export const login: Login = (identifier, password) => {
  const endpoint = getAuthEndpoint(AuthOperation.LOGIN);

  return apiRequest<LoginResponse>(endpoint.path, {
    method: endpoint.method,
    body: { identifier, password },
  });
};
