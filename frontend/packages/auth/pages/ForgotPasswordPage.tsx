import { useEffect, useId, useRef, useState, type FormEvent } from 'react';
import { ArrowLeft, ArrowRight, CheckCircle2, Clock3, Info, LoaderCircle, Mail } from 'lucide-react';
import { Link } from 'react-router-dom';

import { authenticationApi } from '@Business/auth';

import { LoginLayout } from './LoginLayout';

const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
const RESEND_COOLDOWN_MS = 30_000;

/** Keeps the recovery flow predictable without revealing whether an email is registered. */
export const ForgotPasswordPage = () => {
  const id = useId();
  const [email, setEmail] = useState('');
  const [lastRequestedEmail, setLastRequestedEmail] = useState('');
  const [busy, setBusy] = useState(false);
  const [secondsLeft, setSecondsLeft] = useState(0);
  const deadlineRef = useRef(0);
  const requestInFlightRef = useRef(false);
  const [error, setError] = useState('');
  const [emailError, setEmailError] = useState('');

  const countdownActive = secondsLeft > 0;
  const hasRequested = Boolean(lastRequestedEmail);
  const emailChanged = hasRequested && email.trim() !== lastRequestedEmail;

  // The deadline is based on wall-clock time, including when this tab is inactive.
  useEffect(() => {
    if (!countdownActive) return;
    const interval = window.setInterval(() => {
      setSecondsLeft(Math.max(0, Math.ceil((deadlineRef.current - Date.now()) / 1000)));
    }, 250);
    return () => window.clearInterval(interval);
  }, [countdownActive]);

  const submit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (busy || requestInFlightRef.current || Date.now() < deadlineRef.current) return;
    const value = email.trim();
    setError('');
    if (!value) { setEmailError('Please enter your email address.'); return; }
    if (!emailPattern.test(value)) { setEmailError('Enter a valid email address.'); return; }
    setEmailError('');
    requestInFlightRef.current = true;
    setBusy(true);
    try {
      await authenticationApi.forgotPassword(value);
      setLastRequestedEmail(value);
      deadlineRef.current = Date.now() + RESEND_COOLDOWN_MS;
      setSecondsLeft(30);
    } catch {
      setError('Unable to send a reset link right now. Please try again.');
    } finally {
      requestInFlightRef.current = false;
      setBusy(false);
    }
  };

  const cooldown = secondsLeft > 0;
  const noticeSuccess = hasRequested && !emailChanged;
  const buttonLabel = busy
    ? 'Sending link…'
    : cooldown
      ? `Resend in ${secondsLeft}s`
      : hasRequested && !emailChanged
        ? 'Resend reset link'
        : 'Send reset link';

  return <LoginLayout title="Forgot password?" description="We’ll help you get back into your school workspace." eyebrow="ACCOUNT RECOVERY">
    <div className="rdr-auth-flow rdr-auth-flow--forgot">
      <form className="rdr-login-v2__form rdr-auth-flow__form" onSubmit={event => void submit(event)} noValidate>
        <div className="rdr-auth-flow__fields">
          <div className={`rdr-auth-flow__notice ${noticeSuccess ? 'rdr-auth-flow__notice--success' : 'rdr-auth-flow__notice--info'}`}
            role="status" aria-live="polite">
            <span className="rdr-auth-flow__notice-icon" aria-hidden="true">
              {noticeSuccess ? <CheckCircle2 size={20}/> : <Info size={20}/>}
            </span>
            <span className="rdr-auth-flow__notice-copy">
              <strong>{noticeSuccess ? 'Request received!' : emailChanged ? 'Email address updated' : 'Secure password reset'}</strong>
              <small>{noticeSuccess
                ? `If eligible, a password reset link will be sent to ${lastRequestedEmail}.`
                : emailChanged
                  ? 'You can send a link to the corrected address once the countdown ends.'
                  : 'Enter your email to receive a secure password reset link.'}</small>
            </span>
          </div>
          <div className="rdr-login-v2__field rdr-auth-flow__email-field">
            <label htmlFor={id}>Email address</label>
            <div className="rdr-login-v2__input-shell">
              <Mail size={18} aria-hidden="true"/>
              <input id={id} type="email" inputMode="email" autoComplete="email"
                placeholder="you@school.edu" value={email}
                aria-invalid={Boolean(emailError)}
                aria-describedby={emailError ? `${id}-error` : undefined}
                onChange={event => { setEmail(event.target.value); setEmailError(''); setError(''); }} required/>
            </div>
            <small className="rdr-login-v2__field-error rdr-login-v2__field-error--reserved"
              id={`${id}-error`}>{emailError || '\u00a0'}</small>
          </div>
        </div>
        <div className="rdr-auth-flow__actions">
          <div className="rdr-auth-flow__error-slot" aria-live="polite">
            {error && <p className="rdr-login-v2__error" role="alert">{error}</p>}
          </div>
          <button type="submit" className="rdr-login-v2__submit" disabled={busy || cooldown}
            aria-label={buttonLabel}>
            <span>{buttonLabel}</span>
            {busy ? <LoaderCircle className="rdr-login-v2__spinner" size={19}/>
              : cooldown ? <Clock3 size={19}/> : <ArrowRight size={19}/>}
          </button>
          <div className="rdr-auth-flow__secondary-slot">
            <Link className="rdr-auth-flow__back" to="/login"><ArrowLeft size={16}/> Back to sign in</Link>
          </div>
        </div>
      </form>
    </div>
  </LoginLayout>;
};
