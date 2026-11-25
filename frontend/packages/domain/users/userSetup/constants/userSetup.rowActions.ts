import { UserStatus } from '@Business/domain/users';
import type { Account } from '@Business/domain/users';
import type { TableRowAction } from '@Components/components/Table';
import { Eye, Pencil, Mail, UserRoundX } from '@Components/icons';

import { UserAction } from './userSetup.actions';

export const USER_SETUP_ROW_ACTIONS: readonly TableRowAction<Account>[] = Object.freeze([
  { id: UserAction.VIEW, label: __('View account'), icon: Eye },
  { id: UserAction.EDIT, label: __('Edit account'), icon: Pencil },
  { id: UserAction.RESEND, label: __('Resend invitation'), icon: Mail, isVisible: (row) => row.status === UserStatus.PENDING_SETUP },
  { id: UserAction.ACTIVATE, label: __('Activate'), isVisible: (row) => row.status === UserStatus.SUSPENDED },
  { id: UserAction.SUSPEND, label: __('Suspend'), isVisible: (row) => row.status === UserStatus.ACTIVE, destructive: true },
  { id: UserAction.DEACTIVATE, label: __('Deactivate'), icon: UserRoundX, isVisible: (row) => row.status !== UserStatus.DEACTIVATED, destructive: true },
]);
