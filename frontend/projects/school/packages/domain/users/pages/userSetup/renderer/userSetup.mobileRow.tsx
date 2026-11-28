import type { Account } from '@Business/domain/users';
import { UserMobileRow } from '../components/userMobileRow';

export const renderUserMobileRow = (row: Account) => <UserMobileRow row={row} />;
