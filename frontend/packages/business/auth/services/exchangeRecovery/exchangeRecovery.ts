import { apiRequest } from '@Business/api';

import { AuthOperation } from '../constants/auth.operations';
import { getAuthEndpoint } from '../helpers/auth.getEndpoint';
import { getAuthProofHeader } from '../helpers/auth.getProofHeader';

import type { ExchangeRecovery } from './types/exchangeRecovery';
import type { GrantResult } from './types/exchangeRecovery.response';

/** Single-purpose, typed Identity operation. API owns its session cookies. */
export const exchangeRecovery: ExchangeRecovery = (proof) => {
  const endpoint = getAuthEndpoint(AuthOperation.EXCHANGE_RECOVERY);

  return apiRequest<GrantResult>(endpoint.path, {
    method: endpoint.method,
    headers: getAuthProofHeader('RECOVERY', proof),
  });
};
