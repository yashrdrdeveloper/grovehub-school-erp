import { apiRequest } from '@Business/api';

import { AuthOperation } from '../constants/auth.operations';
import { getAuthEndpoint } from '../helpers/auth.getEndpoint';
import { getAuthProofHeader } from '../helpers/auth.getProofHeader';

import type { ValidateInvitation } from './types/validateInvitation';
import type { GrantResult } from './types/validateInvitation.response';

/** Single-purpose, typed Identity operation. API owns its session cookies. */
export const validateInvitation: ValidateInvitation = (proof) => {
  const endpoint = getAuthEndpoint(AuthOperation.VALIDATE_INVITATION);

  return apiRequest<GrantResult>(endpoint.path, {
    method: endpoint.method,
    headers: getAuthProofHeader('INVITATION', proof),
  });
};
