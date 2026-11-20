import { useId, useState, type FormEvent } from 'react';
import { ArrowLeft, ArrowRight, Mail, MailCheck } from 'lucide-react';
import { Link } from 'react-router-dom';

import { authenticationApi } from '@Business/auth';

import { LoginLayout } from './LoginLayout';

const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

export const ForgotPasswordPage = () => {
  const id = useId();
  const [email, setEmail] = useState('');
  const [accepted, setAccepted] = useState(false);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState('');
  const [emailError, setEmailError] = useState('');

  const submit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (busy || accepted) return;
    const value = email.trim();
    setError('');
    if (!value) { setEmailError('Enter your email address.'); return; }
    if (!emailPattern.test(value)) { setEmailError('Enter a valid email address.'); return; }
    setEmailError('');
    setBusy(true);
    try {
      await authenticationApi.forgotPassword(value);
      setAccepted(true);
    } catch {
      setError('Unable to process this request. Try again in a moment.');
    } finally {
      setBusy(false);
    }
  };

  return <LoginLayout title="Forgot your password?" description="Enter your school email and we’ll send you a secure reset link." eyebrow="ACCOUNT RECOVERY">
    <div className="rdr-auth-flow">
      <form className="rdr-login-v2__form rdr-auth-flow__form" onSubmit={event => void submit(event)} noValidate>
        <div className="rdr-auth-flow__fields">
          <div className="rdr-login-v2__field">
            <label htmlFor={id}>Email address</label>
            <div className="rdr-login-v2__input-shell">
              <Mail size={18} aria-hidden="true"/>
              <input id={id} type="email" inputMode="email" autoComplete="email"
                placeholder="you@school.edu" value={email} disabled={accepted}
                aria-invalid={Boolean(emailError)}
                aria-describedby={emailError ? `${id}-error` : undefined}
                onChange={event => {setEmail(event.target.value);setEmailError('');setError('');}} required/>
            </div>
            <small className="rdr-login-v2__field-error rdr-login-v2__field-error--reserved" id={`${id}-error`}>{emailError || '\u00a0'}</small>
          </div>
          <div className="rdr-auth-flow__message" role="status" aria-live="polite">
            {accepted ? <><MailCheck size={19}/> If this account is eligible, recovery instructions are on their way.</> : 'We’ll send instructions if this email belongs to an eligible account.'}
          </div>
        </div>
        <div className="rdr-auth-flow__actions">
          <div className="rdr-login-v2__feedback-slot" aria-live="polite">
            {error && <p className="rdr-login-v2__error" role="alert">{error}</p>}
          </div>
          <button type="submit" className="rdr-login-v2__submit" disabled={busy || accepted}>
            <span>{accepted ? 'Reset link requested' : busy ? 'Sending…' : 'Send reset link'}</span>
            {accepted ? <MailCheck size={19}/> : <ArrowRight size={19}/>}
          </button>
          <Link className="rdr-auth-flow__back" to="/login"><ArrowLeft size={16}/> Back to sign in</Link>
        </div>
      </form>
    </div>
  </LoginLayout>;
};
