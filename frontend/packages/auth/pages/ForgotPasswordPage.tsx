import { useEffect, useId, useRef, useState, type FormEvent } from 'react';
import { ArrowLeft, ArrowRight, CheckCircle2, Clock3, Info, LoaderCircle, Mail, RotateCw } from 'lucide-react';
import { Link, useNavigate } from 'react-router-dom';

import { authenticationApi } from '@Business/auth';

import { LoginLayout } from './LoginLayout';

const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
const RESEND_COOLDOWN_MS = 30_000;

/** Requests a recovery email without disclosing whether an account exists. */
export const ForgotPasswordPage = () => {
  const id = useId();
  const navigate = useNavigate();
  const [email, setEmail] = useState('');
  const [accepted, setAccepted] = useState(false);
  const [busy, setBusy] = useState(false);
  const [secondsLeft, setSecondsLeft] = useState(0);
  const deadlineRef = useRef(0);
  const requestInFlightRef = useRef(false);
  const [error, setError] = useState('');
  const [emailError, setEmailError] = useState('');

  const countdownActive = secondsLeft > 0;
  // Calculate from the wall-clock deadline so background tabs do not slow the countdown.
  useEffect(() => {
    if (!accepted || !countdownActive) return;
    const interval = window.setInterval(() => {
      setSecondsLeft(Math.max(0, Math.ceil((deadlineRef.current - Date.now()) / 1000)));
    }, 250);
    return () => window.clearInterval(interval);
  }, [accepted, countdownActive]);

  const requestLink = async () => {
    if (requestInFlightRef.current || (accepted && Date.now() < deadlineRef.current)) return;
    const value = email.trim();
    setError('');
    if (!value) { setEmailError('Please enter your email address.'); return; }
    if (!emailPattern.test(value)) { setEmailError('Enter a valid email address.'); return; }
    setEmailError('');
    requestInFlightRef.current = true;
    setBusy(true);
    try {
      await authenticationApi.forgotPassword(value);
      setAccepted(true);
      deadlineRef.current = Date.now() + RESEND_COOLDOWN_MS;
      setSecondsLeft(30);
    } catch {
      setError('Unable to request a reset link right now. Please try again.');
    } finally {
      requestInFlightRef.current = false;
      setBusy(false);
    }
  };

  const submit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (accepted) { navigate('/login'); return; }
    void requestLink();
  };

  const canResend = accepted && secondsLeft === 0 && !busy;

  return <LoginLayout title="Forgot password?" description="We'll help you get back into your school workspace." eyebrow="ACCOUNT RECOVERY">
    <div className="rdr-auth-flow rdr-auth-flow--forgot">
      <form className="rdr-login-v2__form rdr-auth-flow__form" onSubmit={submit} noValidate>
        <div className="rdr-auth-flow__fields">
          <div className="rdr-login-v2__field rdr-auth-flow__email-field">
            <label htmlFor={id}>Email address</label>
            <div className="rdr-login-v2__input-shell">
              <Mail size={18} aria-hidden="true"/>
              <input id={id} type="email" inputMode="email" autoComplete="email"
                placeholder="you@school.edu" value={email} readOnly={accepted}
                aria-invalid={Boolean(emailError)}
                aria-describedby={emailError ? `${id}-error` : undefined}
                onChange={event => { setEmail(event.target.value); setEmailError(''); setError(''); }} required/>
            </div>
            <small className="rdr-login-v2__field-error rdr-login-v2__field-error--reserved"
              id={`${id}-error`}>{emailError || '\u00a0'}</small>
          </div>
          <div className={`rdr-auth-flow__notice ${accepted ? 'rdr-auth-flow__notice--success' : 'rdr-auth-flow__notice--info'}`}
            role="status" aria-live="polite">
            <span className="rdr-auth-flow__notice-icon" aria-hidden="true">
              {accepted ? <CheckCircle2 size={20}/> : <Info size={20}/>}
            </span>
            <span className="rdr-auth-flow__notice-copy">
              <strong>{accepted ? 'You’re one step closer!' : 'A secure link, just for you'}</strong>
              <small>{accepted
                ? 'If the address is eligible, a reset link is on its way. Check your inbox and spam folder.'
                : 'We’ll email reset instructions if the address belongs to an eligible account.'}</small>
            </span>
          </div>
          <p className="rdr-auth-flow__hint">
            {accepted ? 'You can request another link after the short cooldown.' : 'For security, we never disclose whether an email is registered.'}
          </p>
        </div>
        <div className="rdr-auth-flow__actions">
          <div className="rdr-auth-flow__error-slot" aria-live="polite">
            {error && <p className="rdr-login-v2__error" role="alert">{error}</p>}
          </div>
          <button type="submit" className="rdr-login-v2__submit" disabled={busy}>
            <span>{accepted ? 'Back to sign in' : busy ? 'Sending link…' : 'Send reset link'}</span>
            {busy ? <LoaderCircle className="rdr-login-v2__spinner" size={19}/> : accepted ? <ArrowLeft size={19}/> : <ArrowRight size={19}/>}
          </button>
          <div className="rdr-auth-flow__secondary-slot">
            {accepted
              ? <button type="button" className="rdr-auth-flow__back rdr-auth-flow__resend"
                  onClick={() => void requestLink()} disabled={!canResend}
                  aria-label={secondsLeft > 0 ? `Resend reset link in ${secondsLeft} seconds` : 'Resend reset link'}>
                  {secondsLeft > 0 ? <Clock3 size={16}/> : <RotateCw size={16}/>}
                  {secondsLeft > 0 ? `Resend link in ${secondsLeft}s` : 'Resend reset link'}
                </button>
              : <Link className="rdr-auth-flow__back" to="/login"><ArrowLeft size={16}/> Back to sign in</Link>}
          </div>
        </div>
      </form>
    </div>
  </LoginLayout>;
};
