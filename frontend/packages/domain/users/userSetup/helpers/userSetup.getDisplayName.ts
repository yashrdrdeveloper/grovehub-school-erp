import type { Account } from '@Business/domain/users';

export const getUserDisplayName = (account: Account): string => {
  if (!account.profile) return account.email;
  return [account.profile.firstName, account.profile.middleName, account.profile.lastName]
    .filter(Boolean)
    .join(' ');
};
