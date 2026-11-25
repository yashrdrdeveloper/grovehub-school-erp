import { getFormattedDate } from '@Base/localization/getFormattedDate';
import { Column } from '@Components/components/Table';
import type { Account } from '@Business/domain/users';

import { renderUserWithEmail } from '../renderer/userSetup.userWithEmail';
import { renderUserStatus } from '../renderer/userSetup.status';
import { UserField } from './userSetup.fields';

export const USER_SETUP_COLUMNS = Object.freeze([
  new Column<Account>()
    .setId(UserField.EMAIL)
    .setHeader(__('User / Email'))
    .setAccessor('email')
    .setRenderer(renderUserWithEmail)
    .setWidth(340)
    .setMinWidth(220)
    .setMaxWidth(4000)
    .setFillRemainingWidth(true)
    .setHideable(false)
    .build(),
  new Column<Account>()
    .setId(UserField.STATUS)
    .setHeader(__('Status'))
    .setAccessor('status')
    .setRenderer(renderUserStatus)
    .setWidth(160)
    .setMinWidth(120)
    .setSortable(true)
    .build(),
  new Column<Account>()
    .setId(UserField.CREATED_AT)
    .setHeader(__('Created at'))
    .setAccessor('createdAt')
    .setRenderer((account) => getFormattedDate(account.createdAt))
    .setWidth(220)
    .setMinWidth(150)
    .setSortable(true)
    .build(),
  new Column<Account>()
    .setId(UserField.UPDATED_AT)
    .setHeader(__('Updated at'))
    .setAccessor('updatedAt')
    .setRenderer((account) => getFormattedDate(account.updatedAt))
    .setWidth(220)
    .setMinWidth(150)
    .setSortable(true)
    .build(),
]);
