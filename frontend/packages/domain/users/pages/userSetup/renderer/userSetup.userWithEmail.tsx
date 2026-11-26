import type { Account } from '@Business/domain/users';
import { UserWithEmail } from '../components/userWithEmail';

export const renderUserWithEmail = (rowData: Account) => <UserWithEmail rowData={rowData} />;
