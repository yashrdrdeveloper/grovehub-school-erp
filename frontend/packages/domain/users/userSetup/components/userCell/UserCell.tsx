import { Avatar } from '@Components/components/Avatar';
import { Box } from '@Components/components/Box';

import { getUserDisplayName } from '../../helpers/userSetup.getDisplayName';
import styles from './UserCell.module.css';
import type { UserCellProps } from './types/userCell.props';

export const UserCell = ({ account }: UserCellProps) => (
  <Box className={styles.user}>
    <Avatar name={getUserDisplayName(account)} />
    <Box className={styles.details}>
      <Box className={styles.name}>{getUserDisplayName(account)}</Box>
      {account.profile && <Box className={styles.email}>{account.email}</Box>}
    </Box>
  </Box>
);
