import { useState, type FormEvent } from 'react';
import { Button } from '@Components/components/Button';
import { Input } from '@Components/components/Input';
import { ShieldCheck } from 'lucide-react';
import { __ } from '@Base/i18n';
import { request, errorText } from '../../../../business/api/client';
import type { Operator } from '../../../../business/api/types';

export function LoginPage({ onLogin }: { onLogin: (operator: Operator) => void }) {
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState('');
  const submit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (busy) return;
    const form = event.currentTarget;
    const data = new FormData(form);
    setBusy(true); setError('');
    try {
      const operator = await request<Operator>('/internal/public/auth/login', 'POST', { email: data.get('email'), password: data.get('password') });
      form.reset(); onLogin(operator);
    } catch (failure) { setError(errorText(failure)); }
    finally { setBusy(false); }
  };
  return <main className="login-layout" data-testid="internal-login-page">
    <aside className="login-brand"><ShieldCheck size={42} /><span className="eyebrow">{__('GROOVEHUB / INTERNAL')}</span><h1>{__('Operations, with accountability.')}</h1><p>{__('A dedicated workspace for approved platform operations, access management and audit history.')}</p></aside>
    <section className="login-panel"><div className="login-card"><span className="eyebrow">{__('OPERATOR WORKSPACE')}</span><h2>{__('Welcome back.')}</h2><p className="muted">{__('Sign in with your Internal Tool operator account.')}</p>
      <form className="form-stack" onSubmit={submit}>
        <Input name="email" label={__('Email address')} type="email" autoComplete="username" required dataTestId="internal-login-email" />
        <Input name="password" label={__('Password')} type="password" autoComplete="current-password" required dataTestId="internal-login-password" />
        {error && <p role="alert" className="error">{error}</p>}
        <Button type="submit" fullWidth loading={busy} loadingText={__('Signing in…')} dataTestId="internal-login-submit">{__('Sign in')}</Button>
      </form><p className="muted small">{__('School accounts cannot sign in to this workspace.')}</p>
    </div></section>
  </main>;
}
