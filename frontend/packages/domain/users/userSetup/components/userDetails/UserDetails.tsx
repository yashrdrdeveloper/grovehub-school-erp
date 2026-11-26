import { _filter } from '@Base/utilities/filter';
import { _map } from '@Base/utilities/map';
import { UserRound } from '@Components/icons';
import { Box } from '@Components/components/Box';
import { Button, BUTTON_VARIANT } from '@Components/components/Button';
import { ButtonWithIcon } from '@Components/components/ButtonWithIcon';
import { Sheet } from '@Components/components/Sheet';

import { UserAction } from '../../constants/userSetup.actions';
import { USER_SETUP_ROW_ACTIONS } from '../../constants/userSetup.rowActions';
import { USER_SETUP_TEST_IDS } from '../../constants/userSetup.testIds';
import { renderCreatedAt, renderUpdatedAt } from '../../renderer/userSetup.timestamps';
import { renderUserStatus } from '../../renderer/userSetup.status';
import styles from '../../UserSetup.module.css';
import type { UserDetailsProps } from './types/userDetails.props';

export const UserDetails = ({ controls, detailRow }: UserDetailsProps) => (
  <Sheet
    open
    title={__('Account details')}
    onOpenChange={(open) => {
      if (!open) controls.closeDialog();
    }}
    dataTestId={USER_SETUP_TEST_IDS.DETAILS}
  >
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
        <Button
          onClick={() =>
            controls.actions.ROW_ACTION?.({ id: UserAction.EDIT, row: detailRow })
          }
        >
          {__('Edit user')}
        </Button>
        {_map(
          _filter(
            USER_SETUP_ROW_ACTIONS,
            (action) => !action.isVisible || action.isVisible(detailRow),
          ),
          (action) => (
            <ButtonWithIcon
              key={action.id}
              icon={action.icon ?? UserRound}
              variant={BUTTON_VARIANT.OUTLINE}
              disabled={action.disabled}
              onClick={() => controls.actions.ROW_ACTION?.({ id: action.id, row: detailRow })}
            >
              {action.label}
            </ButtonWithIcon>
          ),
        )}
      </Box>
    </Box>
  </Sheet>
);
