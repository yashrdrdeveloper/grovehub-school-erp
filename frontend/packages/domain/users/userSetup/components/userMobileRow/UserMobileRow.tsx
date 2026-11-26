import { Avatar } from '@Components/components/Avatar';
import { Box } from '@Components/components/Box';
import { StatusPill } from '@Components/components/StatusPill';
import type { UserMobileRowProps } from './types/userMobileRow.props';
import styles from './UserMobileRow.module.css';

export const UserMobileRow = ({ row }: UserMobileRowProps) => {
  const name = row.username?.trim() || row.email;
  return (
    <Box className={styles.root}>
      <Avatar name={name} />
      <Box className={styles.details}>
        <Box className={styles.name}>{name}</Box>
        <Box className={styles.email}>{row.email}</Box>
        <Box className={styles.status}><StatusPill value={row.status} /></Box>
      </Box>
    </Box>
  );
};
