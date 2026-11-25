import { Box } from '@Components/components/Box';
import { Button, BUTTON_VARIANT } from '@Components/components/Button';
import { Sheet } from '@Components/components/Sheet';
import { Table, TableType, TableDataMode } from '@Components/components/Table';

import { UserForm } from './components/userForm';
import { UserAction } from './constants/userSetup.actions';
import { USER_SETUP_COLUMNS } from './constants/userSetup.columns';
import { USER_SETUP_FILTERS } from './constants/userSetup.filters';
import { USER_SETUP_HEADER } from './constants/userSetup.header';
import { USER_SETUP_ROW_ACTIONS } from './constants/userSetup.rowActions';
import { USER_SETUP_SORT_FIELDS } from './constants/userSetup.sort';
import { USER_SETUP_TEST_IDS } from './constants/userSetup.testIds';
import { getUserSetupRowId } from './helpers/userSetup.getRowId';
import styles from './UserSetup.module.css';
import type { UserSetupProps } from './types/userSetup.props';

export const UserSetup = ({ controls }: UserSetupProps) => {
  const { state, dialog } = controls;
  const getError = () => {
    if (controls.error) return __('Unable to load users. Check your access and try again.');
    return undefined;
  };
  const confirmationTitle = () => {
    if (dialog?.action === UserAction.RESEND) return __('Resend invitation');
    if (dialog?.action === UserAction.ACTIVATE) return __('Activate account');
    if (dialog?.action === UserAction.SUSPEND) return __('Suspend account');
    return __('Deactivate account');
  };
  const showConfirmation = Boolean(dialog?.row && dialog.action !== UserAction.VIEW && dialog.action !== UserAction.EDIT);
  return (
    <Box className={styles.page} data-workspace-fill="true">
      <Table
        id={USER_SETUP_TEST_IDS.TABLE}
        fillHeight
        maxPinnedColumns={3}
        type={TableType.DEFAULT}
        dataMode={TableDataMode.SERVER}
        data={controls.data}
        columns={USER_SETUP_COLUMNS}
        getRowId={getUserSetupRowId}
        headerProps={USER_SETUP_HEADER}
        searchProps={{ value: state.searchText, placeholder: __('Search users...') }}
        filterProps={{ fields: USER_SETUP_FILTERS, value: state.filters, maxFilters: 19 }}
        sortProps={{ fields: USER_SETUP_SORT_FIELDS, value: state.sorting, maxFields: 1 }}
        paginationProps={controls.pagination}
        rowActions={USER_SETUP_ROW_ACTIONS}
        actions={controls.actions}
        loading={controls.loading}
        refreshing={controls.refreshing}
        error={getError()}
      />
      {dialog && (dialog.action === UserAction.CREATE || dialog.action === UserAction.EDIT) && (
        <UserForm row={dialog.row} onClose={controls.closeDialog} onCreate={controls.create} onEdit={controls.edit} submitting={controls.submitting} error={controls.submissionError} />
      )}
      {dialog?.action === UserAction.VIEW && dialog.row && (
        <Sheet open title={__('Account details')} onOpenChange={(open) => { if (!open) controls.closeDialog(); }} dataTestId={USER_SETUP_TEST_IDS.DETAILS}>
          <Box className={styles.details}>
            <Box>{dialog.row.email}</Box>
            <Box>{dialog.row.username ?? '—'}</Box>
            <Box>{dialog.row.status}</Box>
          </Box>
        </Sheet>
      )}
      {showConfirmation && (
        <Sheet open title={confirmationTitle()} onOpenChange={(open) => { if (!open) controls.closeDialog(); }} dataTestId={USER_SETUP_TEST_IDS.CONFIRM}>
          <Box className={styles.details}>
            <Box>{__('Confirm this operation for {{email}}.', { email: dialog?.row?.email ?? '' })}</Box>
            {controls.submissionError && <Box role="alert">{__('Unable to complete the operation.')}</Box>}
            <Box className={styles.actions}>
              <Button variant={BUTTON_VARIANT.OUTLINE} disabled={controls.submitting} onClick={controls.closeDialog}>{__('Cancel')}</Button>
              <Button loading={controls.submitting} onClick={controls.confirm}>{__('Confirm')}</Button>
            </Box>
          </Box>
        </Sheet>
      )}
    </Box>
  );
};
