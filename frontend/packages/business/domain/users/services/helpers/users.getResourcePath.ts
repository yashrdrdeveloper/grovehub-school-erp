import { USERS_ENDPOINT } from '../constants/users.endpoints';

export const getUserResourcePath = (id: string): string =>
  USERS_ENDPOINT.ACCOUNTS + '/' + encodeURIComponent(id);
