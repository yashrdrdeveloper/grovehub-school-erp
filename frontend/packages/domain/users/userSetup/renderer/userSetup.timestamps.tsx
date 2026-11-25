import type { Account } from '@Business/domain/users';
import { UserTimestamp } from '../components/userTimestamp';

export const renderCreatedAt = (row: Account) => <UserTimestamp value={row.createdAt} />;
export const renderUpdatedAt = (row: Account) => <UserTimestamp value={row.updatedAt} />;
