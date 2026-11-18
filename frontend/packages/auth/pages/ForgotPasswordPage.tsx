import { useState, type FormEvent } from 'react';
import { ArrowLeft, MailCheck } from 'lucide-react';
import { Link } from 'react-router-dom';

import { authenticationApi } from '@Business/auth';
import { Button, Input } from '@Components';

import { AuthFrame } from './AuthFrame';

export const ForgotPasswordPage = () => {
  const [email, setEmail] = useState('');
  const [accepted, setAccepted] = useState(false);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState('');
  const [emailError, setEmailError] = useState('');

  const submit = async (event: FormEvent) => {
    event.preventDefault();
    setError('');
    const value = email.trim();
    if (!value) { setEmailError('Enter your email address.'); return; }
    if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(value)) { setEmailError('Enter a valid email address.'); return; }
    setEmailError('');
    setBusy(true);
    try {
      await authenticationApi.forgotPassword(email.trim());
      setAccepted(true);
    } catch {
      setError('Unable to process this request. Try again in a moment.');
    } finally {
      setBusy(false);
    }
  };

  return (
    <AuthFrame title="Forgot password">
                <div className="auth-panel__icon"><MailCheck size={23} /></div>
        <h1>Reset your password</h1>
        <p className="muted">We’ll email a recovery link if this address belongs to an eligible account.</p>
        {accepted ? (
          <div className="notice">If the account is eligible, recovery instructions are on their way.</div>
        ) : (
          <form className="auth-form" noValidate onSubmit={(event) => void submit(event)}>
            <Input id="email" label="Email" placeholder="you@school.edu" type="email" inputMode="email" autoComplete="email" value={email} error={emailError} onChange={(event) => { setEmail(event.target.value); if (emailError) setEmailError(''); }} required />
            {error ? <p className="form-error" role="alert">{error}</p> : null}
            <Button type="submit" loading={busy} className="button--full">Send reset link</Button>
          </form>
        )}
        <Link className="back-link" to="/login"><ArrowLeft size={16} /> Back to sign in</Link>
    </AuthFrame>
  );
};
