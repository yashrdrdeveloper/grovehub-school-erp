import { apiRequest } from '@Business/api';
import { HTTP_METHOD } from '@Business/api/constants/api.method';
import type { PersonNames, SearchRequest, SearchResult } from '@Business/accounts';

import { PERSON_ENDPOINT } from './constants/person.endpoints';
import { getPersonResourcePath } from './helpers/person.getResourcePath';

import type { Person } from '../types/person.entity';

export const personApi = {
  create: (person: PersonNames, contacts: unknown, key: string) =>
    apiRequest<Person>(PERSON_ENDPOINT.PROTECTED, {
      method: HTTP_METHOD.POST,
      body: { ...person, contacts },
      headers: { 'Idempotency-Key': key },
    }),

  get: (gid: string) => apiRequest<Person>(getPersonResourcePath(gid)),

  update: (gid: string, patch: Partial<Person>) =>
    apiRequest<Person>(getPersonResourcePath(gid), {
      method: HTTP_METHOD.PATCH,
      body: patch,
    }),

  search: (request: SearchRequest) =>
    apiRequest<SearchResult<Person>>(`${PERSON_ENDPOINT.PROTECTED}/search`, {
      method: HTTP_METHOD.POST,
      body: request,
    }),

  deactivate: (gid: string) =>
    apiRequest<Person>(`${getPersonResourcePath(gid)}/deactivate`, {
      method: HTTP_METHOD.PATCH,
    }),
};
