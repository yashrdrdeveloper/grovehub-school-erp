import { apiRequest } from '@Business/api';
import type { PersonNames, SearchRequest, SearchResult } from '@Business/accounts';

export interface Person {
  gid: string;
  firstName: string;
  middleName: string | null;
  lastName: string;
  status: string;
  version: number;
}

const ROOT = '/foundation/protected/persons';

export const personApi = {
  create: (person: PersonNames, contacts: unknown, key: string) =>
    apiRequest<Person>(ROOT, { method: 'POST', body: { ...person, contacts }, headers: { 'Idempotency-Key': key } }),
  get: (gid: string) => apiRequest<Person>(`${ROOT}/${encodeURIComponent(gid)}`),
  update: (gid: string, patch: Partial<Person>) =>
    apiRequest<Person>(`${ROOT}/${encodeURIComponent(gid)}`, { method: 'PATCH', body: patch }),
  search: (request: SearchRequest) =>
    apiRequest<SearchResult<Person>>(`${ROOT}/search`, { method: 'POST', body: request }),
  deactivate: (gid: string) =>
    apiRequest<Person>(`${ROOT}/${encodeURIComponent(gid)}/deactivate`, { method: 'PATCH' }),
};
