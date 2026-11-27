import { LoginForm } from './components/LoginForm';
import { useLoginSubmission } from './hooks/useLoginSubmission';
import { useLoginViewport } from './hooks/useLoginViewport';
import { getLoginBrandTemplate } from './renderer/login.formTemplate';
import type { LoginPageProps } from './types/login.props';
import styles from './Login.module.css';

export function LoginPage(props: LoginPageProps) {
  const form = useLoginSubmission(props);
  const mobile = useLoginViewport();
  const Brand = getLoginBrandTemplate(mobile);
  return (
    <main className={styles.layout} data-testid="internal-login-page">
      <Brand />
      <LoginForm {...form} />
    </main>
  );
}
