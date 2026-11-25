import { _map } from '@Base/utilities/map';
import type { Account } from '@Business/domain/users';
import { Table, TableType, TableDataMode, Column } from '@Components/components/Table';
import { FilterFieldType, FilterOperator } from '@Components/components/Filters';
import type { FilterField } from '@Components/components/Filters';
import { StatusPill } from '@Components/components/StatusPill';

import { UserCell } from '../userCell';
import { USER_SETUP_COLUMNS } from '../../constants/userSetup.columns';
import { USER_SETUP_FILTERS, USER_SETUP_QUICK_FILTERS } from '../../constants/userSetup.filters';
import { USER_SETUP_HEADER } from '../../constants/userSetup.header';
import { USER_SETUP_ROW_ACTIONS } from '../../constants/userSetup.rowActions';
import { USER_SETUP_SORT_FIELDS } from '../../constants/userSetup.sort';
import { USER_SETUP_TEST_IDS } from '../../constants/userSetup.testIds';
import { UserField } from '../../constants/userSetup.fields';
import { getUserSetupRowId } from '../../helpers/userSetup.getRowId';
import { getUserRoleLabel } from '../../helpers/userSetup.getRoleLabel';

import type { UserTableProps } from './types/userTable.props';

export const UserTable = ({ controls }: UserTableProps) => {
  const { state, roles } = controls;
  const roleOptions = _map(roles, (role) => ({ value: role.id, label: role.name }));
  const fields: FilterField[] = [...USER_SETUP_FILTERS];
  if (roleOptions.length) {
    fields.push({
      id: UserField.PRIMARY_ROLE,
      label: __('Primary role'),
      type: FilterFieldType.SELECT,
      operators: [FilterOperator.EQUALS, FilterOperator.NOT_EQUALS, FilterOperator.IN, FilterOperator.NOT_IN],
      options: roleOptions,
    });
  }
  const configured = _map(USER_SETUP_COLUMNS, (column) => {
    if (column.id === UserField.NAME) return { ...column, renderer: (account: Account) => <UserCell account={account} /> };
    if (column.id === UserField.STATUS) return { ...column, renderer: (account: Account) => <StatusPill value={account.status} /> };
    return column;
  });
  const columns = [
    ...configured.slice(0, 1),
    new Column<Account>()
      .setId(UserField.PRIMARY_ROLE)
      .setHeader(__('Primary role'))
      .setAccessor((account) => getUserRoleLabel(roles, account.primaryRoleId))
      .setWidth(180)
      .build(),
    new Column<Account>()
      .setId(UserField.SECONDARY_ROLES)
      .setHeader(__('Secondary roles'))
      .setAccessor((account) => _map(account.secondaryRoleIds ?? [], (id) => getUserRoleLabel(roles, id)).join(', ') || '—')
      .setWidth(220)
      .build(),
    ...configured.slice(1),
  ];
  const getError = () => {
    if (controls.error) return __('Unable to load users. Check your access and try again.');
    return undefined;
  };
  return (
    <Table
      id={USER_SETUP_TEST_IDS.TABLE}
      type={TableType.DEFAULT}
      dataMode={TableDataMode.SERVER}
      data={controls.data}
      columns={columns}
      getRowId={getUserSetupRowId}
      headerProps={USER_SETUP_HEADER}
      searchProps={{ value: state.searchText, placeholder: __('Search users...') }}
      filterProps={{ fields, value: state.filters, maxFilters: 19 }}
      sortProps={{ fields: USER_SETUP_SORT_FIELDS, value: state.sorting, maxFields: 1 }}
      quickFilterProps={{ options: USER_SETUP_QUICK_FILTERS, value: state.quickFilterIds }}
      paginationProps={controls.pagination}
      rowActions={USER_SETUP_ROW_ACTIONS}
      actions={controls.actions}
      preferencesProps={{ storageKey: controls.preferencesKey }}
      loading={controls.loading}
      refreshing={controls.refreshing}
      error={getError()}
    />
  );
};
