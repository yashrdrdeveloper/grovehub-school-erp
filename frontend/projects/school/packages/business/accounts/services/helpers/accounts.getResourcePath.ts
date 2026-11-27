import { ACCOUNTS_ENDPOINT } from '../constants/accounts.endpoints';

import type { GetAccountResourcePath } from './types/accounts.getResourcePath';

/** URL-encode untrusted account IDs exactly once. */
export const getAccountResourcePath: GetAccountResourcePath = (id) =>
  `${ACCOUNTS_ENDPOINT.PROTECTED}/${encodeURIComponent(id)}`;
