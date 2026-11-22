import { apiRequest } from '@Business/api';

import { AuthOperation } from '../constants/auth.operations';
import { getAuthEndpoint } from '../helpers/auth.getEndpoint';
import { getAuthProofHeader } from '../helpers/auth.getProofHeader';

import type { ValidateRecovery } from './types/validateRecovery';
import type { GrantResult } from './types/validateRecovery.response';

/** Single-purpose, typed Identity operation. API owns its session cookies. */
export const validateRecovery: ValidateRecovery = (proof) => {
  const endpoint = getAuthEndpoint(AuthOperation.VALIDATE_RECOVERY);

  return apiRequest<GrantResult>(endpoint.path, {
    method: endpoint.method,
    headers: getAuthProofHeader('RECOVERY', proof),
  });
};
