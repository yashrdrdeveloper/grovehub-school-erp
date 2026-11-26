import { _filter } from '@Base/utilities/filter';
import { _map } from '@Base/utilities/map';
import { UserRound } from '@Components/icons';
import { Box } from '@Components/components/Box';
import { Button, BUTTON_VARIANT } from '@Components/components/Button';
import { ButtonWithIcon } from '@Components/components/ButtonWithIcon';
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
import { renderUserMobileRow } from './renderer/userSetup.mobileRow';
import { renderCreatedAt, renderUpdatedAt } from './renderer/userSetup.timestamps';
import { renderUserStatus } from './renderer/userSetup.status';
import { getUserSetupRowId } from './helpers/userSetup.getRowId';
import styles from './UserSetup.module.css';
import type { UserSetupProps } from './types/userSetup.props';

export const UserSetup = ({ controls }: UserSetupProps) => {
  const { state, dialog } = controls;
  const detailRow = dialog?.row;
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
        totalRecords={controls.totalRecords}
        columns={USER_SETUP_COLUMNS}
        getRowId={getUserSetupRowId}
        headerProps={USER_SETUP_HEADER}
        mobileRowProps={{ render: renderUserMobileRow, rowActions: [], onClick: (row) => controls.actions.ROW_ACTION?.({ id: UserAction.VIEW, row }) }}
        mobileInfiniteProps={controls.mobileInfiniteProps}
        searchProps={{ value: state.searchText, placeholder: __('Search by username or email...'), mobilePlaceholder: __('Search users'), fields: [__('Username'), __('Email address')] }}
        filterProps={{ fields: USER_SETUP_FILTERS, value: state.filters, maxFilters: USER_SETUP_FILTERS.length }}
        sortProps={{ fields: USER_SETUP_SORT_FIELDS, value: state.sorting, maxFields: USER_SETUP_SORT_FIELDS.length }}
        paginationProps={controls.pagination}
        rowActions={USER_SETUP_ROW_ACTIONS}
        actions={controls.actions}
        loading={controls.loading}
        refreshing={controls.refreshing}
        error={getError()}
      />
      {dialog && (dialog.action === UserAction.CREATE || dialog.action === UserAction.EDIT) && (
        <UserForm row={dialog.row} onClose={controls.closeDialog} onCreate={controls.create} onEdit={controls.edit} onAccountAction={(id, row) => controls.actions.ROW_ACTION?.({ id, row })} submitting={controls.submitting} error={controls.submissionError} />
      )}
      {dialog?.action === UserAction.VIEW && detailRow && (
        <Sheet open title={__('Account details')} onOpenChange={(open) => { if (!open) controls.closeDialog(); }} dataTestId={USER_SETUP_TEST_IDS.DETAILS}>
          <Box className={styles.details}>
            <Box className={styles.detailIdentity}>
              <Box className={styles.detailName}>{detailRow.username ?? '—'}</Box>
              <Box className={styles.detailEmail}>{detailRow.email}</Box>
              <Box>{renderUserStatus(detailRow)}</Box>
            </Box>
            <Box className={styles.detailField}>
              <Box>{__('Created at')}</Box>
              {renderCreatedAt(detailRow)}
            </Box>
            <Box className={styles.detailField}>
              <Box>{__('Updated at')}</Box>
              {renderUpdatedAt(detailRow)}
            </Box>
            <Box className={styles.detailActions}>
              <Button onClick={() => controls.actions.ROW_ACTION?.({ id: UserAction.EDIT, row: detailRow })}>{__('Edit user')}</Button>
              {_map(_filter(USER_SETUP_ROW_ACTIONS, (action) => !action.isVisible || action.isVisible(detailRow)), (action) => (
                <ButtonWithIcon key={action.id} icon={action.icon ?? UserRound} variant={BUTTON_VARIANT.OUTLINE} disabled={action.disabled} onClick={() => controls.actions.ROW_ACTION?.({ id: action.id, row: detailRow })}>{action.label}</ButtonWithIcon>
              ))}
            </Box>
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
