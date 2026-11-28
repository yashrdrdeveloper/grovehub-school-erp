import { Button } from '@Components/components/Button';
import { Input } from '@Components/components/Input';
import { PasswordField } from '@Components/components/PasswordField';

import { LOGIN_COPY } from '../constants/login.copy';
import type { LoginFormProps } from '../types/login.props';
import styles from '../Login.module.css';

export function LoginForm({ busy, error, onSubmit }: LoginFormProps) {
  return (
    <section className={styles.panel}>
      <div className={styles.card}>
        <span className={styles.eyebrow}>{LOGIN_COPY.workspace}</span>
        <h1>{LOGIN_COPY.title}</h1>
        <p className={styles.muted}>{LOGIN_COPY.description}</p>
        <form className={styles.form} onSubmit={onSubmit} data-testid="internal-login-form">
          <Input name="email" label={LOGIN_COPY.email} type="email" autoComplete="username" required dataTestId="internal-login-email" />
          <PasswordField name="password" label={LOGIN_COPY.password} autoComplete="current-password" required dataTestId="internal-login-password" />
          {error && <p role="alert" className={styles.error}>{error}</p>}
          <Button type="submit" fullWidth loading={busy} loadingText={LOGIN_COPY.busyButton} dataTestId="internal-login-submit">{LOGIN_COPY.button}</Button>
        </form>
        <p className={styles.note}>{LOGIN_COPY.note}</p>
      </div>
    </section>
  );
}
