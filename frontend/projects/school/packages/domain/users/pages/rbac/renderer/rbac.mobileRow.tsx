import type { Role } from '@Business/domain/users/rbac/types/contracts';
import styles from '../page.module.css';

export const renderRbacMobileRow = (row: Role) => <div className={styles.row}><strong>{row.name}</strong><span>{row.status}</span><span>{row.description}</span></div>;
