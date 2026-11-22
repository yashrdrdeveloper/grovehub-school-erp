import { apiRequest } from '@Business/api';
import { HTTP_METHOD } from '@Business/api/constants/api.method';
import type { Account, SetupResult } from '@Business/auth';

import { ACCOUNTS_ENDPOINT } from './constants/accounts.endpoints';
import { getAccountResourcePath } from './helpers/accounts.getResourcePath';

import type {
  AccountSearchResult,
  CreateAccount,
  SearchRequest,
} from '../types';
import type {
  CreateAccountResponse,
  AccountSetupStatus,
  ResendInvitationResponse,
} from '../types/account.responses';

/** Account transport; business entities and browser navigation are separate concerns. */
export const accountsApi = {
  create: (input: CreateAccount, idempotencyKey: string) =>
    apiRequest<CreateAccountResponse>(ACCOUNTS_ENDPOINT.PROTECTED, {
      method: HTTP_METHOD.POST,
      body: input,
      headers: { 'Idempotency-Key': idempotencyKey },
    }),

  search: (request: SearchRequest) =>
    apiRequest<AccountSearchResult>(`${ACCOUNTS_ENDPOINT.PROTECTED}/search`, {
      method: HTTP_METHOD.POST,
      body: request,
    }),

  get: (accountId: string) => apiRequest<Account>(getAccountResourcePath(accountId)),

  update: (accountId: string, username: string | null) =>
    apiRequest<Account>(getAccountResourcePath(accountId), {
      method: HTTP_METHOD.PATCH,
      body: { username },
    }),

  suspend: (accountId: string) =>
    apiRequest<Account>(`${getAccountResourcePath(accountId)}/suspend`, {
      method: HTTP_METHOD.PATCH,
    }),

  activate: (accountId: string) =>
    apiRequest<Account>(`${getAccountResourcePath(accountId)}/activate`, {
      method: HTTP_METHOD.PATCH,
    }),

  deactivate: (accountId: string) =>
    apiRequest<Account>(`${getAccountResourcePath(accountId)}/deactivate`, {
      method: HTTP_METHOD.PATCH,
    }),

  resendInvitation: (accountId: string) =>
    apiRequest<ResendInvitationResponse>(
      `${getAccountResourcePath(accountId)}/resend-invitation`,
      { method: HTTP_METHOD.POST },
    ),

  setupStatus: () =>
    apiRequest<AccountSetupStatus>(ACCOUNTS_ENDPOINT.PUBLIC_SETUP_STATUS),

  setup: (password: string) =>
    apiRequest<SetupResult>(ACCOUNTS_ENDPOINT.PUBLIC_SETUP, {
      method: HTTP_METHOD.POST,
      body: { password },
    }),
};
