import { apiRequest } from '@Business/api';
import { HTTP_METHOD } from '@Business/api/constants/api.method';
import type { SearchRequest, SearchResult } from '@Business/accounts/types';
import type { Role } from '../types/contracts';

export const rbacApi = {
  search: (request: SearchRequest) =>
    apiRequest<SearchResult<Role>>('/identity/protected/roles/search', {
      method: HTTP_METHOD.POST,
      body: request,
    }),
  permissions: (id: string) =>
    apiRequest<{ permissionCodes: string[]; version: number }>(
      `/identity/protected/roles/${encodeURIComponent(id)}/permissions`,
    ),
};
