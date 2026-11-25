import { getFormattedDate } from '@Base/localization/getFormattedDate';
import { Column } from '@Components/components/Table';
import type { Account } from '@Business/domain/users';

import { UserField } from './userSetup.fields';

export const USER_SETUP_COLUMNS = Object.freeze([
  new Column<Account>()
    .setId(UserField.EMAIL)
    .setHeader(__('Email address'))
    .setAccessor('email')
    .setHideable(false)
    .build(),
  new Column<Account>()
    .setId(UserField.USERNAME)
    .setHeader(__('Username'))
    .setAccessor('username')
    .setSortable(true)
    .build(),
  new Column<Account>()
    .setId(UserField.STATUS)
    .setHeader(__('Status'))
    .setAccessor('status')
    .setSortable(true)
    .build(),
  new Column<Account>()
    .setId(UserField.CREATED_AT)
    .setHeader(__('Created at'))
    .setAccessor('createdAt')
    .setRenderer((account) => getFormattedDate(account.createdAt))
    .setWidth(190)
    .setSortable(true)
    .build(),
  new Column<Account>()
    .setId(UserField.UPDATED_AT)
    .setHeader(__('Updated at'))
    .setAccessor('updatedAt')
    .setRenderer((account) => getFormattedDate(account.updatedAt))
    .setWidth(190)
    .setSortable(true)
    .build(),
]);
