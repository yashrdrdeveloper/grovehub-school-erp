import { apiRequest } from '@Business/api';

import { AuthOperation } from '../constants/auth.operations';
import { getAuthEndpoint } from '../helpers/auth.getEndpoint';
import { getAuthProofHeader } from '../helpers/auth.getProofHeader';

import type { ExchangeInvitation } from './types/exchangeInvitation';
import type { GrantResult } from './types/exchangeInvitation.response';

/** Single-purpose, typed Identity operation. API owns its session cookies. */
export const exchangeInvitation: ExchangeInvitation = (proof) => {
  const endpoint = getAuthEndpoint(AuthOperation.EXCHANGE_INVITATION);

  return apiRequest<GrantResult>(endpoint.path, {
    method: endpoint.method,
    headers: getAuthProofHeader('INVITATION', proof),
  });
};
