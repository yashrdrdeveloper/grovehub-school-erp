import { apiRequest } from '@Business/api';
import { HTTP_METHOD } from '@Business/api/constants/api.method';
import type { SearchRequest, SearchResult } from '@Business/accounts/types';
import type { AuditActivity } from '../types/contracts';

export const auditApi = {
  search: (request: SearchRequest & { scope?: 'ME' }) => apiRequest<SearchResult<AuditActivity>>('/audit/protected/search', { method: HTTP_METHOD.POST, body: request }),
};
