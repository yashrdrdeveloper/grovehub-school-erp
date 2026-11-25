import type { Account } from '@Business/domain/users';
import { StatusPill } from '@Components/components/StatusPill';

export const renderUserStatus = (account: Account) => <StatusPill value={account.status} />;
