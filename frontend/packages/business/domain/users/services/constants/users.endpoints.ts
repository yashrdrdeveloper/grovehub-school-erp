import { ACCOUNTS_ENDPOINT } from '@Business/accounts/services/constants/accounts.endpoints';

export const USERS_ENDPOINT = Object.freeze({
  ACCOUNTS: ACCOUNTS_ENDPOINT.PROTECTED,
  SEARCH: ACCOUNTS_ENDPOINT.PROTECTED + '/search',
  ROLES: ACCOUNTS_ENDPOINT.PROTECTED + '/roles',
});
