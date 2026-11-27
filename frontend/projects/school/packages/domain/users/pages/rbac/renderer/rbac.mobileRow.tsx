import { Box } from '@Components/components/Box';
import { Typography, TYPOGRAPHY_TAG } from '@Components/components/Typography';
import type { Role } from '@Business/domain/users/rbac/types/contracts';
import styles from '../page.module.css';

export const renderRbacMobileRow = (row: Role) => (
  <Box className={styles.row}>
    <Typography as={TYPOGRAPHY_TAG.STRONG}>{row.name}</Typography>
    <Typography as={TYPOGRAPHY_TAG.SPAN}>{row.status}</Typography>
    <Typography as={TYPOGRAPHY_TAG.SPAN}>{row.description}</Typography>
  </Box>
);
