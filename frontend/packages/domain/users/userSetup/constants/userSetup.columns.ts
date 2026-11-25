import { Column, TablePin } from '@Components/components/Table';
import type { Account } from '@Business/domain/users';

import { getUserDisplayName } from '../helpers/userSetup.getDisplayName';
import { UserField } from './userSetup.fields';

export const USER_SETUP_COLUMNS = Object.freeze([
  new Column<Account>()
    .setId(UserField.NAME)
    .setHeader(__('User'))
    .setAccessor(getUserDisplayName)
    .setWidth(280)
    .setPin(TablePin.LEFT)
    .setHideable(false)
    .build(),
  new Column<Account>()
    .setId(UserField.USERNAME)
    .setHeader(__('Username'))
    .setAccessor('username')
    .setMobileVisible(false)
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
    .setMobileVisible(false)
    .setSortable(true)
    .build(),
  new Column<Account>()
    .setId(UserField.UPDATED_AT)
    .setHeader(__('Updated at'))
    .setAccessor('updatedAt')
    .setMobileVisible(false)
    .setVisible(false)
    .setSortable(true)
    .build(),
]);
