import { ShieldCheck } from 'lucide-react';
import { LOGIN_COPY } from '../../constants/login.copy';
import styles from '../../Login.module.css';

export function DesktopLoginBrand() {
  return (
    <aside className={styles.brand} data-testid="internal-login-brand">
      <ShieldCheck size={42} aria-hidden="true" />
      <span className={styles.brandEyebrow}>{LOGIN_COPY.brand}</span>
      <p className={styles.brandTitle}>{LOGIN_COPY.brandTitle}</p>
      <p className={styles.brandDescription}>{LOGIN_COPY.brandDescription}</p>
    </aside>
  );
}
