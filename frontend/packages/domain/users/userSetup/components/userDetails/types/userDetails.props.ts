import type { Account } from '@Business/domain/users';
import type { UserSetupProps } from '../../../types/userSetup.props';

export interface UserDetailsProps extends UserSetupProps {
  readonly detailRow: Account;
}
