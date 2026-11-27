import { accountsApi } from '@Business/accounts';

export const usersApi = {
  search: accountsApi.search,
  create: accountsApi.create,
  get: accountsApi.get,
  update: accountsApi.update,
  activate: accountsApi.activate,
  suspend: accountsApi.suspend,
  deactivate: accountsApi.deactivate,
  resendInvitation: accountsApi.resendInvitation,
};
