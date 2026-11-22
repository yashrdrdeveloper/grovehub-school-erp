import { apiRequest } from '@Business/api';

import { AuthOperation } from '../constants/auth.operations';
import { getAuthEndpoint } from '../helpers/auth.getEndpoint';
import { getAuthProofHeader } from '../helpers/auth.getProofHeader';
import { login } from './login';

import type {
  Account,
  GrantResult,
  RefreshResult,
  SessionSummary,
  StatusResult,
} from '../types';

const endpoint = getAuthEndpoint;

/** An operation-indexed, exhaustively typed facade without widened any signatures. */
export const AUTHENTICATION_OPERATIONS = {
  [AuthOperation.LOGIN]: login,

  [AuthOperation.REFRESH]: () =>
    apiRequest<RefreshResult>(endpoint(AuthOperation.REFRESH).path, {
      method: endpoint(AuthOperation.REFRESH).method,
    }),

  [AuthOperation.FORGOT_PASSWORD]: (email: string) =>
    apiRequest<{ accepted: boolean; delivery: string }>(
      endpoint(AuthOperation.FORGOT_PASSWORD).path,
      {
        method: endpoint(AuthOperation.FORGOT_PASSWORD).method,
        body: { email },
      },
    ),

  [AuthOperation.VALIDATE_RECOVERY]: (proof: string) =>
    apiRequest<GrantResult>(endpoint(AuthOperation.VALIDATE_RECOVERY).path, {
      method: endpoint(AuthOperation.VALIDATE_RECOVERY).method,
      headers: getAuthProofHeader('RECOVERY', proof),
    }),

  [AuthOperation.VALIDATE_INVITATION]: (proof: string) =>
    apiRequest<GrantResult>(endpoint(AuthOperation.VALIDATE_INVITATION).path, {
      method: endpoint(AuthOperation.VALIDATE_INVITATION).method,
      headers: getAuthProofHeader('INVITATION', proof),
    }),

  [AuthOperation.EXCHANGE_RECOVERY]: (proof: string) =>
    apiRequest<GrantResult>(endpoint(AuthOperation.EXCHANGE_RECOVERY).path, {
      method: endpoint(AuthOperation.EXCHANGE_RECOVERY).method,
      headers: getAuthProofHeader('RECOVERY', proof),
    }),

  [AuthOperation.EXCHANGE_INVITATION]: (proof: string) =>
    apiRequest<GrantResult>(endpoint(AuthOperation.EXCHANGE_INVITATION).path, {
      method: endpoint(AuthOperation.EXCHANGE_INVITATION).method,
      headers: getAuthProofHeader('INVITATION', proof),
    }),

  [AuthOperation.RESET_GRANT_STATUS]: () =>
    apiRequest<GrantResult>(endpoint(AuthOperation.RESET_GRANT_STATUS).path, {
      method: endpoint(AuthOperation.RESET_GRANT_STATUS).method,
    }),

  [AuthOperation.RESET_PASSWORD]: (newPassword: string) =>
    apiRequest<StatusResult>(endpoint(AuthOperation.RESET_PASSWORD).path, {
      method: endpoint(AuthOperation.RESET_PASSWORD).method,
      body: { newPassword },
    }),

  [AuthOperation.ME]: () =>
    apiRequest<Account>(endpoint(AuthOperation.ME).path, {
      method: endpoint(AuthOperation.ME).method,
    }),

  [AuthOperation.LOGOUT]: () =>
    apiRequest<StatusResult>(endpoint(AuthOperation.LOGOUT).path, {
      method: endpoint(AuthOperation.LOGOUT).method,
    }),

  [AuthOperation.LOGOUT_ALL]: () =>
    apiRequest<{ status: string; revokedCount: number }>(
      endpoint(AuthOperation.LOGOUT_ALL).path,
      { method: endpoint(AuthOperation.LOGOUT_ALL).method },
    ),

  [AuthOperation.CHANGE_PASSWORD]: (newPassword: string) =>
    apiRequest<StatusResult>(endpoint(AuthOperation.CHANGE_PASSWORD).path, {
      method: endpoint(AuthOperation.CHANGE_PASSWORD).method,
      body: { newPassword },
    }),

  [AuthOperation.REAUTHENTICATE]: (password: string) =>
    apiRequest<StatusResult>(endpoint(AuthOperation.REAUTHENTICATE).path, {
      method: endpoint(AuthOperation.REAUTHENTICATE).method,
      body: { password },
    }),

  [AuthOperation.SESSIONS]: () =>
    apiRequest<{ items: SessionSummary[] }>(endpoint(AuthOperation.SESSIONS).path, {
      method: endpoint(AuthOperation.SESSIONS).method,
    }),

  [AuthOperation.REVOKE_SESSION]: (sessionId: string) =>
    apiRequest<{ sessionId: string; status: string }>(
      endpoint(AuthOperation.REVOKE_SESSION).path + encodeURIComponent(sessionId),
      { method: endpoint(AuthOperation.REVOKE_SESSION).method },
    ),
} as const;
