import { apiRequest } from '@Business/api';

import { AuthOperation } from '../../constants/auth.operations';
import { getAuthEndpoint } from '../../helpers/auth.getEndpoint';

import type { ResetGrantStatus } from './types/resetGrantStatus';
import type { GrantResult } from './types/resetGrantStatus.response';

/** Single-purpose, typed Identity operation. API owns its session cookies. */
export const resetGrantStatus: ResetGrantStatus = () => {
  const endpoint = getAuthEndpoint(AuthOperation.RESET_GRANT_STATUS);

  return apiRequest<GrantResult>(endpoint.path, {
    method: endpoint.method,
  });
};
