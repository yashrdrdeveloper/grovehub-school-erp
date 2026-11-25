import { FilterFieldType, FilterOperator } from '@Components/components/Filters';
import type { FilterField } from '@Components/components/Filters';
import { UserStatus } from '@Business/domain/users';

import { UserField } from './userSetup.fields';

export const USER_SETUP_STATUS_OPTIONS = Object.freeze([
  { value: UserStatus.PENDING_SETUP, label: __('Pending setup') },
  { value: UserStatus.ACTIVE, label: __('Active') },
  { value: UserStatus.SUSPENDED, label: __('Suspended') },
  { value: UserStatus.DEACTIVATED, label: __('Deactivated') },
]);

export const USER_SETUP_FILTERS: readonly FilterField[] = Object.freeze([
  {
    id: UserField.EMAIL,
    label: __('Email address'),
    type: FilterFieldType.TEXT,
    operators: [FilterOperator.EQUALS, FilterOperator.NOT_EQUALS, FilterOperator.CONTAINS, FilterOperator.STARTS_WITH],
  },
  {
    id: UserField.USERNAME,
    label: __('Username'),
    type: FilterFieldType.TEXT,
    operators: [FilterOperator.EQUALS, FilterOperator.NOT_EQUALS, FilterOperator.CONTAINS, FilterOperator.STARTS_WITH],
  },
  {
    id: UserField.STATUS,
    label: __('Status'),
    type: FilterFieldType.SELECT,
    operators: [FilterOperator.EQUALS, FilterOperator.NOT_EQUALS, FilterOperator.IN, FilterOperator.NOT_IN],
    options: USER_SETUP_STATUS_OPTIONS,
  },
]);

export const USER_SETUP_QUICK_FILTERS = Object.freeze([
  { id: 'ALL', label: __('All') },
  { id: UserStatus.ACTIVE, label: __('Active') },
  { id: UserStatus.PENDING_SETUP, label: __('Pending setup') },
  { id: UserStatus.SUSPENDED, label: __('Suspended') },
  { id: UserStatus.DEACTIVATED, label: __('Deactivated') },
]);
