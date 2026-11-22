import { apiRequest } from '@Business/api';

import { AuthOperation } from '../../constants/auth.operations';
import { getAuthEndpoint } from '../../helpers/auth.getEndpoint';

import type { Refresh } from './types/refresh';
import type { RefreshResult } from './types/refresh.response';

/** Single-purpose, typed Identity operation. API owns its session cookies. */
export const refresh: Refresh = () => {
  const endpoint = getAuthEndpoint(AuthOperation.REFRESH);

  return apiRequest<RefreshResult>(endpoint.path, {
    method: endpoint.method,
  });
};
