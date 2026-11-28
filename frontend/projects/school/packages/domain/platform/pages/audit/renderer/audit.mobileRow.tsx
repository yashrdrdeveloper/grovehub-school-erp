import { Box } from '@Components/components/Box';
import { Typography, TYPOGRAPHY_TAG } from '@Components/components/Typography';
import type { AuditActivity } from '@Business/domain/platform/audit/types/contracts';
import styles from '../page.module.css';

export const renderAuditMobileRow = (row: AuditActivity) => (
  <Box className={styles.row}>
    <Typography as={TYPOGRAPHY_TAG.STRONG}>{row.operation}</Typography>
    <Typography as={TYPOGRAPHY_TAG.SPAN}>{row.source}</Typography>
    <Typography as={TYPOGRAPHY_TAG.SPAN}>{row.assetType}</Typography>
  </Box>
);
