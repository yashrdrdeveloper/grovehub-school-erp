import { ShieldCheck } from 'lucide-react';
import { LOGIN_COPY } from '../../constants/login.copy';
import styles from '../../Login.module.css';

export function MobileLoginBrand() {
  return <div className={styles.mobileBrand} data-testid="internal-login-mobile-brand">
    <ShieldCheck size={28} aria-hidden="true" />
    <span>{LOGIN_COPY.brand}</span>
  </div>;
}
