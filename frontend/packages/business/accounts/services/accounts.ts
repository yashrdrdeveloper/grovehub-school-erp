import { apiRequest } from '@Business/api';
import type { Account, SetupResult } from '@Business/auth';

import type { AccountSearchResult, CreateAccount, SearchRequest } from '../types';

const PROTECTED = '/identity/protected/accounts';

export const accountsApi = {
  create: (input: CreateAccount, idempotencyKey: string) =>
    apiRequest<{
      accountId: string;
      personGid: string;
      email: string;
      status: string;
      invitation: { status: string };
    }>(PROTECTED, {
      method: 'POST',
      body: input,
      headers: { 'Idempotency-Key': idempotencyKey },
    }),
  search: (request: SearchRequest) =>
    apiRequest<AccountSearchResult>(`${PROTECTED}/search`, {
      method: 'POST',
      body: request,
    }),
  get: (accountId: string) =>
    apiRequest<Account>(`${PROTECTED}/${encodeURIComponent(accountId)}`),
  update: (accountId: string, username: string | null) =>
    apiRequest<Account>(`${PROTECTED}/${encodeURIComponent(accountId)}`, {
      method: 'PATCH',
      body: { username },
    }),
  suspend: (accountId: string) =>
    apiRequest<Account>(`${PROTECTED}/${encodeURIComponent(accountId)}/suspend`, { method: 'PATCH' }),
  activate: (accountId: string) =>
    apiRequest<Account>(`${PROTECTED}/${encodeURIComponent(accountId)}/activate`, { method: 'PATCH' }),
  deactivate: (accountId: string) =>
    apiRequest<Account>(`${PROTECTED}/${encodeURIComponent(accountId)}/deactivate`, { method: 'PATCH' }),
  resendInvitation: (accountId: string) =>
    apiRequest<{ accountId: string; invitationStatus: string }>(
      `${PROTECTED}/${encodeURIComponent(accountId)}/resend-invitation`,
      { method: 'POST' },
    ),
  setup: (username: string, password: string) =>
    apiRequest<SetupResult>('/identity/public/accounts/setup', {
      method: 'POST',
      body: { username, password },
    }),
};
