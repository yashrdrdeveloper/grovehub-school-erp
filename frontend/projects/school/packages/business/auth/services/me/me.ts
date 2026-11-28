import { apiRequest } from '@Business/api';

import { AuthOperation } from '../constants/auth.operations';
import { getAuthEndpoint } from '../helpers/auth.getEndpoint';

import type { Me } from './types/me';
import type { Account } from './types/me.response';

/** Single-purpose, typed Identity operation. API owns its session cookies. */
export const me: Me = () => {
  const endpoint = getAuthEndpoint(AuthOperation.ME);

  return apiRequest<Account>(endpoint.path, {
    method: endpoint.method,
  });
};
