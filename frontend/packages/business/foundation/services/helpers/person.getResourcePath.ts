import { PERSON_ENDPOINT } from '../constants/person.endpoints';

import type { GetPersonResourcePath } from './types/person.getResourcePath';

/** Consistent, encoded person URL for all Foundation transport methods. */
export const getPersonResourcePath: GetPersonResourcePath = (gid) =>
  `${PERSON_ENDPOINT.PROTECTED}/${encodeURIComponent(gid)}`;
