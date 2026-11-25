import { UserRound } from '@Components/icons';
import { Avatar } from '@Components/components/Avatar';
import { Box } from '@Components/components/Box';
import { Icon } from '@Components/components/Icon';

import type { UserWithEmailProps } from './types/userWithEmail.props';
import styles from './UserWithEmail.module.css';

export const UserWithEmail = ({ rowData }: UserWithEmailProps) => {
  const username = rowData.username?.trim();
  const name = username || __('Unknown user');
  const getFallback = () => {
    if (!username) return <Icon icon={UserRound} size={18} />;
    return undefined;
  };
  return (
    <Box className={styles.root}>
      <Avatar name={name} fallback={getFallback()} />
      <Box className={styles.details}>
        <Box className={styles.name} title={name}>{name}</Box>
        <Box className={styles.email} title={rowData.email}>{rowData.email}</Box>
      </Box>
    </Box>
  );
};
