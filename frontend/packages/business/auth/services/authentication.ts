import { apiRequest } from '@Business/api';

import type {
  Account,
  GrantResult,
  LoginResult,
  RefreshResult,
  SessionSummary,
  StatusResult,
} from '../types';

const PUBLIC = '/identity/public/auth';
const PROTECTED = '/identity/protected/auth';

export const authenticationApi = {
  login: (identifier: string, password: string) =>
    apiRequest<LoginResult>(`${PUBLIC}/login`, {
      method: 'POST',
      body: { identifier, password },
    }),
  refresh: () =>
    apiRequest<RefreshResult>(`${PUBLIC}/refresh`, { method: 'POST' }),
  forgotPassword: (email: string) =>
    apiRequest<{ accepted: boolean; delivery: string }>(`${PUBLIC}/forgot-password`, {
      method: 'POST',
      body: { email },
    }),
  exchangeRecovery: (proof: string) =>
    apiRequest<GrantResult>(`${PUBLIC}/recovery/exchange`, {
      method: 'POST',
      headers: { Authorization: `Recovery ${proof}` },
    }),
  exchangeInvitation: (proof: string) =>
    apiRequest<GrantResult>(`${PUBLIC}/invitations/exchange`, {
      method: 'POST',
      headers: { Authorization: `Invitation ${proof}` },
    }),
  resetPassword: (newPassword: string) =>
    apiRequest<StatusResult>(`${PUBLIC}/reset-password`, {
      method: 'POST',
      body: { newPassword },
    }),
  me: () => apiRequest<Account>(`${PROTECTED}/me`),
  logout: () => apiRequest<StatusResult>(`${PROTECTED}/logout`, { method: 'POST' }),
  logoutAll: () => apiRequest<{ status: string; revokedCount: number }>(`${PROTECTED}/logout-all`, { method: 'POST' }),
  changePassword: (newPassword: string) =>
    apiRequest<StatusResult>(`${PROTECTED}/change-password`, {
      method: 'POST',
      body: { newPassword },
    }),
  reauthenticate: (password: string) =>
    apiRequest<StatusResult>(`${PROTECTED}/reauthenticate`, {
      method: 'POST',
      body: { password },
    }),
  sessions: () => apiRequest<{ items: SessionSummary[] }>(`${PROTECTED}/sessions`),
  revokeSession: (sessionId: string) =>
    apiRequest<{ sessionId: string; status: string }>(
      `${PROTECTED}/sessions/${encodeURIComponent(sessionId)}`,
      { method: 'DELETE' },
    ),
};
