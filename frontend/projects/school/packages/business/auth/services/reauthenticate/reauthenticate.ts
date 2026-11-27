import { apiRequest } from '@Business/api';

import { AuthOperation } from '../constants/auth.operations';
import { getAuthEndpoint } from '../helpers/auth.getEndpoint';

import type { Reauthenticate } from './types/reauthenticate';
import type { StatusResult } from './types/reauthenticate.response';

/** Single-purpose, typed Identity operation. API owns its session cookies. */
export const reauthenticate: Reauthenticate = (password) => {
  const endpoint = getAuthEndpoint(AuthOperation.REAUTHENTICATE);

  return apiRequest<StatusResult>(endpoint.path, {
    method: endpoint.method,
    body: { password },
  });
};
