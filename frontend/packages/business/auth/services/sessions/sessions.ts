import { apiRequest } from '@Business/api';

import { AuthOperation } from '../constants/auth.operations';
import { getAuthEndpoint } from '../helpers/auth.getEndpoint';

import type { Sessions } from './types/sessions';
import type { SessionsResponse } from './types/sessions.response';

/** Single-purpose, typed Identity operation. API owns its session cookies. */
export const sessions: Sessions = () => {
  const endpoint = getAuthEndpoint(AuthOperation.SESSIONS);

  return apiRequest<SessionsResponse>(endpoint.path, {
    method: endpoint.method,
  });
};
