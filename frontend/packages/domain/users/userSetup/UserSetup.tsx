import { Box } from '@Components/components/Box';
import { Button, BUTTON_VARIANT } from '@Components/components/Button';
import { Sheet } from '@Components/components/Sheet';
import { UserTable } from './components/userTable';

import { UserForm } from './components/userForm';
import { UserAction } from './constants/userSetup.actions';
import { USER_SETUP_TEST_IDS } from './constants/userSetup.testIds';
import { getUserRoleLabel } from './helpers/userSetup.getRoleLabel';
import styles from './UserSetup.module.css';
import type { UserSetupProps } from './types/userSetup.props';

export const UserSetup = ({ controls }: UserSetupProps) => {
  const { dialog } = controls;
  const confirmationTitle = () => {
    if (dialog?.action === UserAction.RESEND) return __('Resend invitation');
    if (dialog?.action === UserAction.ACTIVATE) return __('Activate account');
    if (dialog?.action === UserAction.SUSPEND) return __('Suspend account');
    return __('Deactivate account');
  };
  const showConfirmation = Boolean(dialog?.row && dialog.action !== UserAction.VIEW && dialog.action !== UserAction.EDIT);
  return (
    <Box className={styles.page} data-workspace-fill="true">
      <UserTable controls={controls} />
      {dialog && (dialog.action === UserAction.CREATE || dialog.action === UserAction.EDIT) && (
        <UserForm
          key={dialog.row?.accountId ?? UserAction.CREATE}
          row={dialog.row}
          roles={controls.roles}
          rolesLoading={controls.rolesLoading}
          rolesError={controls.rolesError}
          retryRoles={controls.retryRoles}
          profile={controls.profile}
          profileLoading={controls.profileLoading}
          profileError={controls.profileError}
          retryProfile={controls.retryProfile}
          onClose={controls.closeDialog}
          onCreate={controls.create}
          onEdit={controls.edit}
          onEditProfile={controls.editProfile}
          submitting={controls.submitting}
          error={controls.submissionError}
          profileSubmitting={controls.profileSubmitting}
          profileSubmissionError={controls.profileSubmissionError}
          profileSaved={controls.profileSaved}
        />
      )}
      {dialog?.action === UserAction.VIEW && dialog.row && (
        <Sheet open title={__('Account details')} onOpenChange={(open) => { if (!open) controls.closeDialog(); }} dataTestId={USER_SETUP_TEST_IDS.DETAILS}>
          <Box className={styles.details}>
            <Box>{dialog.row.email}</Box>
            <Box>{dialog.row.username ?? '—'}</Box>
            <Box>{dialog.row.status}</Box>
            <Box>{__('Primary role')}: {getUserRoleLabel(controls.roles, dialog.row.primaryRoleId)}</Box>
            <Box>{__('Secondary roles')}: {(dialog.row.secondaryRoleIds ?? []).map((id) => getUserRoleLabel(controls.roles, id)).join(', ') || '—'}</Box>
            {controls.profile && <Box>{[controls.profile.firstName, controls.profile.middleName, controls.profile.lastName].filter(Boolean).join(' ')}</Box>}
            {controls.profile?.phoneNumber && <Box>{controls.profile.phoneNumber}</Box>}
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
