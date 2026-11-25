import { accountsApi } from '@Business/accounts';
import type { SearchRequest, SearchResult } from '@Business/accounts';
import { apiRequest } from '@Business/api';
import { HTTP_METHOD } from '@Business/api/constants/api.method';

import { USERS_ENDPOINT } from './constants/users.endpoints';
import { getUserResourcePath } from './helpers/users.getResourcePath';
import type { Account, CreateUser, UpdateUser, UpdateUserProfile, UserProfile } from '../types/users.contracts';
import type { UserRole } from '../types/users.roles';

export const usersApi = {
  search: (request: SearchRequest, signal?: AbortSignal) =>
    apiRequest<SearchResult<Account>>(USERS_ENDPOINT.SEARCH, {
      method: HTTP_METHOD.POST,
      body: request,
      signal,
    }),
  create: (input: CreateUser, key: string) => accountsApi.create(input, key),
  get: (id: string) => apiRequest<Account>(getUserResourcePath(id)),
  roles: (signal?: AbortSignal) => apiRequest<readonly UserRole[]>(USERS_ENDPOINT.ROLES, { signal }),
  update: (id: string, input: UpdateUser) =>
    apiRequest<Account>(getUserResourcePath(id), { method: HTTP_METHOD.PATCH, body: input }),
  profile: (id: string, signal?: AbortSignal) =>
    apiRequest<UserProfile>(getUserResourcePath(id) + '/profile', { signal }),
  updateProfile: (id: string, input: UpdateUserProfile) =>
    apiRequest<UserProfile>(getUserResourcePath(id) + '/profile', { method: HTTP_METHOD.PATCH, body: input }),
  activate: accountsApi.activate,
  suspend: accountsApi.suspend,
  deactivate: accountsApi.deactivate,
  resendInvitation: accountsApi.resendInvitation,
};
