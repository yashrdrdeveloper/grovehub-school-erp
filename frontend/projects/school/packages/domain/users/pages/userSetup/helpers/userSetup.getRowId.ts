import type { Account } from '@Business/domain/users';

export const getUserSetupRowId = (account: Account): string => account.accountId;
