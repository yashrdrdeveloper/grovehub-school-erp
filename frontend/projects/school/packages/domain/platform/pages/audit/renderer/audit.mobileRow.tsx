import type { AuditActivity } from '@Business/domain/platform/audit/types/contracts';
import styles from '../page.module.css';

export const renderAuditMobileRow = (row: AuditActivity) => <div className={styles.row}><strong>{row.operation}</strong><span>{row.source}</span><span>{row.assetType}</span></div>;
