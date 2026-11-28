import { SortDirection } from '@Components/components/Sort';

import { UserField } from './userSetup.fields';

export const USER_SETUP_SORT_FIELDS = Object.freeze([
  { id: UserField.USERNAME, label: __('Username') },
  { id: UserField.STATUS, label: __('Status') },
  { id: UserField.CREATED_AT, label: __('Created at') },
  { id: UserField.UPDATED_AT, label: __('Updated at') },
]);

export const USER_SETUP_DEFAULT_SORT = Object.freeze([
  { field: UserField.CREATED_AT, direction: SortDirection.DESC },
]);
