import { useEffect, useId, useRef, useState, type FormEvent } from 'react';
import { ArrowLeft, ArrowRight, CheckCircle2, Eye, EyeOff, KeyRound, LockKeyhole, ShieldCheck, UserRound } from 'lucide-react';
import { Link } from 'react-router-dom';

import { authenticationApi } from '@Business/auth';
import { accountsApi } from '@Business/accounts';

import { LoginLayout } from './LoginLayout';

type Purpose = 'recovery' | 'invitation';

export const RecoveryPage = ({ purpose }: { purpose: Purpose }) => {
  const id = useId();
  const proofRef = useRef<string | null>(new URLSearchParams(window.location.hash.slice(1)).get('proof'));
  const verificationRef = useRef<Promise<unknown> | null>(null);
  const [verified, setVerified] = useState(false);
  const [password, setPassword] = useState('');
  const [confirmation, setConfirmation] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [username, setUsername] = useState('');
  const [complete, setComplete] = useState(false);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState('');
  const [fieldErrors, setFieldErrors] = useState<{ username?: string; password?: string; confirmation?: string }>({});

  useEffect(() => {
    // Keep the proof in the URL fragment only until exchange; do not send it in a request URL.
    const proof = proofRef.current;
    window.history.replaceState(null, '', window.location.pathname);
    if (!proof) {
      setError('This link is missing its one-time proof. Request a fresh link.');
      return;
    }
    let active = true;
    if (!verificationRef.current) {
      verificationRef.current = purpose === 'recovery'
        ? authenticationApi.exchangeRecovery(proof)
        : authenticationApi.exchangeInvitation(proof);
    }
    void verificationRef.current.then(
      () => { if (active) setVerified(true); },
      () => { if (active) setError('This link is expired or already used. Request a fresh one.'); },
    );
    return () => { active = false; };
  }, [purpose]);

  const submit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (busy || !verified || complete) return;
    setError('');
    const next: typeof fieldErrors = {};
    if (purpose === 'invitation' && !username.trim()) next.username = 'Choose a username.';
    if (!password) next.password = 'Enter a new password.';
    else if (password.length < 12) next.password = 'Use at least 12 characters.';
    if (!confirmation) next.confirmation = 'Confirm your new password.';
    else if (password !== confirmation) next.confirmation = 'Passwords do not match.';
    setFieldErrors(next);
    if (Object.keys(next).length) return;
    setBusy(true);
    try {
      if (purpose === 'recovery') await authenticationApi.resetPassword(password);
      else await accountsApi.setup(username.trim(), password);
      setComplete(true);
      setPassword('');
      setConfirmation('');
    } catch (failure) {
      setError(failure instanceof Error ? failure.message : 'Unable to complete setup. Try again.');
    } finally {
      setBusy(false);
    }
  };

  const settingUp = purpose === 'invitation';
  const title = settingUp ? 'Set your password.' : 'Create a new password.';
  const description = settingUp
    ? 'Complete your school invitation and secure your account.'
    : 'Choose a strong password to get back into your school workspace.';

  return <LoginLayout title={title} description={description} eyebrow={settingUp ? 'WELCOME TO RDR SCHOOL' : 'SECURE ACCOUNT RECOVERY'}>
    <div className="rdr-auth-flow">
      <form className="rdr-login-v2__form rdr-auth-flow__form" onSubmit={event => void submit(event)} noValidate>
        <div className="rdr-auth-flow__fields">
          {settingUp && <div className="rdr-login-v2__field">
            <label htmlFor={`${id}-username`}>Username</label>
            <div className="rdr-login-v2__input-shell">
              <UserRound size={18} aria-hidden="true"/>
              <input id={`${id}-username`} autoComplete="username" placeholder="Choose a username"
                disabled={!verified || complete} value={username}
                aria-invalid={Boolean(fieldErrors.username)}
                onChange={event => {setUsername(event.target.value);setFieldErrors(current => ({...current,username:undefined}));}} required/>
            </div>
            <small className="rdr-login-v2__field-error rdr-login-v2__field-error--reserved">{fieldErrors.username || '\u00a0'}</small>
          </div>}
          <div className="rdr-login-v2__field">
            <label htmlFor={`${id}-password`}>New password</label>
            <div className="rdr-login-v2__input-shell">
              <LockKeyhole size={18} aria-hidden="true"/>
              <input id={`${id}-password`} type={showPassword ? 'text' : 'password'}
                autoComplete="new-password" placeholder="At least 12 characters"
                disabled={!verified || complete} value={password} minLength={12}
                aria-invalid={Boolean(fieldErrors.password)}
                onChange={event => {setPassword(event.target.value);setFieldErrors(current => ({...current,password:undefined}));}} required/>
              <button type="button" className="rdr-login-v2__visibility"
                aria-label={showPassword ? 'Hide passwords' : 'Show passwords'}
                onClick={() => setShowPassword(value => !value)}>
                {showPassword ? <EyeOff size={19}/> : <Eye size={19}/>}
              </button>
            </div>
            <small className="rdr-login-v2__field-error rdr-login-v2__field-error--reserved">{fieldErrors.password || '\u00a0'}</small>
          </div>
          <div className="rdr-login-v2__field">
            <label htmlFor={`${id}-confirmation`}>Confirm new password</label>
            <div className="rdr-login-v2__input-shell">
              <KeyRound size={18} aria-hidden="true"/>
              <input id={`${id}-confirmation`} type={showPassword ? 'text' : 'password'}
                autoComplete="new-password" placeholder="Re-enter your password"
                disabled={!verified || complete} value={confirmation}
                aria-invalid={Boolean(fieldErrors.confirmation)}
                onChange={event => {setConfirmation(event.target.value);setFieldErrors(current => ({...current,confirmation:undefined}));}} required/>
            </div>
            <small className="rdr-login-v2__field-error rdr-login-v2__field-error--reserved">{fieldErrors.confirmation || '\u00a0'}</small>
          </div>
          <div className="rdr-auth-flow__message" role="status" aria-live="polite">
            {complete
              ? <><CheckCircle2 size={19}/> {settingUp ? 'Your school account is ready.' : 'Your password has been reset.'}</>
              : !verified
                ? <><ShieldCheck size={19}/> {error || 'Verifying your secure link…'}</>
                : <><ShieldCheck size={19}/> Use 12 or more characters for your new password.</>}
          </div>
        </div>
        <div className="rdr-auth-flow__actions">
          <div className="rdr-login-v2__feedback-slot" aria-live="polite">
            {verified && error && <p className="rdr-login-v2__error" role="alert">{error}</p>}
          </div>
          {complete
            ? <Link className="rdr-login-v2__submit rdr-auth-flow__submit-link" to="/login"><span>Continue to sign in</span><ArrowRight size={19}/></Link>
            : <button className="rdr-login-v2__submit" type="submit" disabled={busy || !verified}>
                <span>{busy ? 'Saving securely…' : settingUp ? 'Activate account' : 'Reset password'}</span>
                <ArrowRight size={19}/>
              </button>}
          <Link className="rdr-auth-flow__back" to="/login"><ArrowLeft size={16}/> Back to sign in</Link>
        </div>
      </form>
    </div>
  </LoginLayout>;
};
