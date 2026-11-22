import { apiRequest } from '@Business/api';

import { AuthOperation } from '../../constants/auth.operations';
import { getAuthEndpoint } from '../../helpers/auth.getEndpoint';

import type { RevokeSession } from './types/revokeSession';
import type { RevokeSessionResponse } from './types/revokeSession.response';

/** Single-purpose, typed Identity operation. API owns its session cookies. */
export const revokeSession: RevokeSession = (sessionId) => {
  const endpoint = getAuthEndpoint(AuthOperation.REVOKE_SESSION);

  return apiRequest<RevokeSessionResponse>(endpoint.path + encodeURIComponent(sessionId), {
    method: endpoint.method,
  });
};
