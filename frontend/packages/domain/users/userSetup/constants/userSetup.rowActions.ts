import { UserStatus } from '@Business/domain/users';
import type { Account } from '@Business/domain/users';
import type { TableRowAction } from '@Components/components/Table';
import { ShieldOff, UserRoundX } from '@Components/icons';

import { UserAction } from './userSetup.actions';

export const USER_SETUP_ROW_ACTIONS: readonly TableRowAction<Account>[] = Object.freeze([
  { id: UserAction.SUSPEND, label: __('Suspend'), icon: ShieldOff, isVisible: (row) => row.status === UserStatus.ACTIVE },
  { id: UserAction.DEACTIVATE, label: __('Deactivate'), icon: UserRoundX, isVisible: (row) => row.status !== UserStatus.DEACTIVATED },
]);
