import { useEffect, useId, useRef, useState, type FormEvent } from 'react';
import { AlertCircle, ArrowLeft, Clock3, Eye, EyeOff, KeyRound, Link2Off, LoaderCircle, LockKeyhole, PartyPopper, ShieldCheck, Sparkles } from 'lucide-react';
import { Link } from 'react-router-dom';

import { authenticationApi } from '@Business/auth';
import { accountsApi } from '@Business/accounts';
import { ApiError } from '@Business/api';

import { LoginLayout } from './LoginLayout';
import { AuthMobileStatusArt } from './AuthMobileArtwork';
import { clearRecoveryFlow, readRecoveryFlow, saveRecoveryGrant, saveRecoveryProof } from './recoveryFlowSession';

type Purpose = 'recovery' | 'invitation';
type ViewState = 'missing' | 'checking' | 'expired' | 'unavailable' | 'ready' | 'saving' | 'success' | 'attempts-exhausted';
const MAX_SETUP_ATTEMPTS = 2;

const RecoveryArtwork = ({ state, invitation }: { state: ViewState; invitation: boolean }) => {
  const success = state === 'success';
  const checking = state === 'checking';
  const Icon = success ? PartyPopper : checking ? ShieldCheck : state === 'missing' ? Link2Off : state === 'attempts-exhausted' ? AlertCircle : Clock3;
  const mobileStatusTitle = success
    ? invitation ? 'Account ready' : 'Password updated'
    : checking ? 'Checking your link'
    : state === 'missing' ? 'This link looks incomplete'
    : state === 'attempts-exhausted' ? 'Setup unavailable'
    : state === 'unavailable' ? 'Unable to check your link'
    : 'This link has expired';
  const title = success
    ? invitation ? 'Welcome to RDR School!' : 'Password updated!'
    : checking ? 'Checking your secure link'
    : state === 'missing' ? 'This link looks incomplete'
    : state === 'attempts-exhausted' ? 'We couldn’t finish your setup'
    : state === 'unavailable' ? 'We couldn’t check your link'
    : 'This link is no longer valid';
  const message = success
    ? invitation ? 'Your account is ready. A new chapter begins here.' : 'Your password was reset. You can now sign in securely.'
    : checking ? 'We’re verifying your invitation and preparing your workspace.'
    : state === 'missing' ? invitation ? 'Please contact your RDR School administrator for a fresh invitation.' : 'Please request another password reset link.'
    : state === 'attempts-exhausted' ? 'Please contact your RDR School administrator to regenerate your invitation.'
    : state === 'unavailable' ? 'A connection issue prevented verification. Your link may still be valid.'
    : invitation ? 'It may have expired or already been used. Ask your school administrator to regenerate the invitation.' : 'Your reset link has expired or has already been used. Please request a new one.';
  return <div className={`rdr-recovery__art rdr-recovery__art--${success ? 'success' : checking ? 'checking' : 'problem'}`} role="status" aria-live="polite">
    <AuthMobileStatusArt kind={state === 'success' ? 'success' : state === 'checking' ? 'checking' : state === 'expired' ? 'expired' : state === 'unavailable' ? 'unavailable' : state === 'attempts-exhausted' ? 'attempts-exhausted' : 'missing'}/>
    <strong className="rdr-recovery__mobile-status-title">{mobileStatusTitle}</strong>
    <div className="rdr-recovery__art-visual" aria-hidden="true">
      <span className="rdr-recovery__art-orbit rdr-recovery__art-orbit--one"/>
      <span className="rdr-recovery__art-orbit rdr-recovery__art-orbit--two"/>
      <span className="rdr-recovery__art-icon"><Icon size={44} strokeWidth={1.55}/></span>
      {success ? <Sparkles className="rdr-recovery__sparkle" size={25}/> : null}
    </div>
    <strong className="rdr-recovery__art-title">{title}</strong>
    <p className="rdr-recovery__art-message">{message}</p>
  </div>;
};

export const RecoveryPage = ({ purpose }: { purpose: Purpose }) => {
  const invitation = purpose === 'invitation';
  const id = useId();
  // Fragment is read before the first paint and then scrubbed immediately.
  const proofRef = useRef<string | null>(null);
  const flowSequence = useRef(0);
  const attemptRef = useRef(0);
  const submittingRef = useRef(false);
  const [view, setView] = useState<ViewState>('checking');
  const [slow, setSlow] = useState(false);
  const [password, setPassword] = useState('');
  const [confirmation, setConfirmation] = useState('');
  const [visible, setVisible] = useState(false);
  const [error, setError] = useState('');
  const [fieldErrors, setFieldErrors] = useState<{ password?: string; confirmation?: string }>({});
  const [attempts, setAttempts] = useState(0);

  const verify = async (proof: string, sequence: number) => {
    setView('checking');
    try {
      if (invitation) await authenticationApi.validateInvitation(proof);
      else await authenticationApi.validateRecovery(proof);
      if (sequence !== flowSequence.current) return;
      const grant = invitation
        ? await authenticationApi.exchangeInvitation(proof)
        : await authenticationApi.exchangeRecovery(proof);
      if (sequence !== flowSequence.current) return;
      // The proof is single-use. Persist only a marker for the HttpOnly
      // setup/reset grant cookie so a page reload can resume via status.
      saveRecoveryGrant(purpose, grant.expiresIn);
      proofRef.current = null;
      setView('ready');
    } catch (failure) {
      if (sequence !== flowSequence.current) return;
      const expired = failure instanceof ApiError &&
        (failure.code === 'IDENTITY_TOKEN_INVALID' || [401, 403, 404, 410].includes(failure.status));
      if (expired) { clearRecoveryFlow(); proofRef.current = null; }
      setView(expired ? 'expired' : 'unavailable');
    }
  };

  useEffect(() => {
    const handleLink = () => {
      const url = new URL(window.location.href);
      const incoming = new URLSearchParams(url.hash.slice(1)).get('proof')?.trim();
      if (incoming) {
        // A newly pasted fragment always takes precedence over a saved flow.
        // Reset the form and invalidate older in-flight validation results.
        const sequence = ++flowSequence.current;
        saveRecoveryProof(purpose, incoming);
        proofRef.current = incoming;
        attemptRef.current = 0;
        setAttempts(0);
        setPassword('');
        setConfirmation('');
        setError('');
        setFieldErrors({});
        window.history.replaceState(window.history.state, '', url.pathname + url.search);
        void verify(incoming, sequence);
        return;
      }

      // Do not scrub unrelated URL fragments. An absent proof on its own
      // must never contact the backend.
      const saved = readRecoveryFlow(purpose);
      if (!saved) {
        setView('missing');
        return;
      }
      const sequence = ++flowSequence.current;
      if (saved.stage === 'proof') {
        proofRef.current = saved.proof;
        void verify(saved.proof, sequence);
        return;
      }

      // Only an earlier successful one-time exchange may check grant status.
      // Without a saved marker, the page makes zero requests.
      proofRef.current = null;
      setView('checking');
      void (async () => {
        try {
          if (invitation) await accountsApi.setupStatus();
          else await authenticationApi.resetGrantStatus();
          if (sequence === flowSequence.current) setView('ready');
        } catch (failure) {
          if (sequence !== flowSequence.current) return;
          if (failure instanceof ApiError && failure.code === 'IDENTITY_SETUP_ATTEMPTS_EXCEEDED') {
            clearRecoveryFlow();
            setView('attempts-exhausted');
          } else if (failure instanceof ApiError && [401, 403, 404, 410].includes(failure.status)) {
            clearRecoveryFlow();
            setView('expired');
          } else {
            setView('unavailable');
          }
        }
      })();
    };

    handleLink();
    // Changing #proof=... without a page reload must start a fresh flow.
    window.addEventListener('hashchange', handleLink);
    return () => {
      ++flowSequence.current;
      window.removeEventListener('hashchange', handleLink);
    };
  }, [purpose, invitation]);

  useEffect(() => {
    if (view !== 'checking') { setSlow(false); return; }
    const timer = window.setTimeout(() => setSlow(true), 5000);
    return () => window.clearTimeout(timer);
  }, [view]);

  const submit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    if (view !== 'ready' || submittingRef.current) return;
    const next: typeof fieldErrors = {};
    if (!password) next.password = 'Please enter a password.';
    else if (password.length < 12) next.password = 'Use at least 12 characters.';
    if (!confirmation) next.confirmation = 'Please confirm your password.';
    else if (confirmation !== password) next.confirmation = 'Passwords do not match.';
    setFieldErrors(next);
    if (Object.keys(next).length > 0) return; // Never count client-side validation.
    setError('');
    submittingRef.current = true;
    setView('saving');
    try {
      if (invitation) await accountsApi.setup(password);
      else await authenticationApi.resetPassword(password);
      setPassword('');
      setConfirmation('');
      clearRecoveryFlow();
      setView('success');
    } catch (failure) {
      const apiError = failure instanceof ApiError ? failure : null;
      const attemptsExceeded = apiError?.code === 'IDENTITY_SETUP_ATTEMPTS_EXCEEDED';
      const invalid = apiError?.code === 'IDENTITY_TOKEN_INVALID' || [401, 403, 410].includes(apiError?.status ?? 0);
      const nextAttempts = attemptRef.current + 1;
      attemptRef.current = nextAttempts;
      setAttempts(nextAttempts);
      setPassword('');
      setConfirmation('');
      if (attemptsExceeded) { clearRecoveryFlow(); setView('attempts-exhausted'); }
      else if (invalid) { clearRecoveryFlow(); setView('expired'); }
      else if (invitation && nextAttempts >= MAX_SETUP_ATTEMPTS) { clearRecoveryFlow(); setView('attempts-exhausted'); }
      else {
        setError(invitation ? 'Activation failed. You have one more attempt.' : 'Unable to reset your password. Please try again.');
        setView('ready');
      }
    } finally {
      submittingRef.current = false;
    }
  };

  const formVisible = view === 'ready' || view === 'saving';
  const title = 'Set a new password.';
  // Keep the page identity stable throughout validation, errors and success.
  const mobileTitle = invitation ? 'Set password' : 'Reset password';
  return <LoginLayout variant="recovery" title={title} description={invitation
    ? 'Choose a strong password to protect your school account.'
    : 'Choose a strong password to protect your school account.'}
    mobileDescription={formVisible ? 'Choose a strong password to protect your school account.' : ''}
    minimalMobile mobileTitle={mobileTitle}
    eyebrow="SECURE YOUR ACCOUNT">
    <div className="rdr-recovery">
      {formVisible ? <form className="rdr-login-v2__form rdr-recovery__form" onSubmit={event => void submit(event)} noValidate>
        <div className="rdr-recovery__fields">
          <div className="rdr-login-v2__field">
            <label htmlFor={`${id}-password`}>New password</label>
            <div className="rdr-login-v2__input-shell">
              <LockKeyhole size={18} aria-hidden="true"/>
              <input id={`${id}-password`} type={visible ? 'text' : 'password'} autoComplete="new-password"
                placeholder="New password" value={password} minLength={12} disabled={view === 'saving'}
                aria-invalid={Boolean(fieldErrors.password)}
                aria-describedby={fieldErrors.password ? `${id}-password-error` : undefined}
                onChange={event => { setPassword(event.target.value);setFieldErrors(e => ({...e,password:undefined})); }} required/>
              <button type="button" className="rdr-login-v2__visibility" aria-label={visible ? 'Hide passwords' : 'Show passwords'}
                onClick={() => setVisible(!visible)} disabled={view === 'saving'}>
                {visible ? <EyeOff size={19} aria-hidden="true"/> : <Eye size={19} aria-hidden="true"/>}
              </button>
            </div>
            <small id={`${id}-password-error`} className={`rdr-login-v2__field-error rdr-login-v2__field-error--reserved${fieldErrors.password ? "" : " is-empty"}`} aria-live="polite" aria-atomic="true">{fieldErrors.password || '\u00a0'}</small>
          </div>
          <div className="rdr-login-v2__field">
            <label htmlFor={`${id}-confirmation`}>Confirm password</label>
            <div className="rdr-login-v2__input-shell">
              <KeyRound size={18} aria-hidden="true"/>
              <input id={`${id}-confirmation`} type={visible ? 'text' : 'password'} autoComplete="new-password"
                placeholder="Confirm password" value={confirmation} disabled={view === 'saving'}
                aria-invalid={Boolean(fieldErrors.confirmation)}
                aria-describedby={fieldErrors.confirmation ? `${id}-confirmation-error` : undefined}
                onChange={event => { setConfirmation(event.target.value);setFieldErrors(e => ({...e,confirmation:undefined})); }} required/>
              <button type="button" className="rdr-login-v2__visibility"
                aria-label={visible ? 'Hide confirmation password' : 'Show confirmation password'}
                onClick={() => setVisible(current => !current)} disabled={view === 'saving'}>
                {visible ? <EyeOff size={19} aria-hidden="true"/> : <Eye size={19} aria-hidden="true"/>}
              </button>
            </div>
            <small id={`${id}-confirmation-error`} className={`rdr-login-v2__field-error rdr-login-v2__field-error--reserved${fieldErrors.confirmation ? "" : " is-empty"}`} aria-live="polite" aria-atomic="true">{fieldErrors.confirmation || '\u00a0'}</small>
          </div>
          <p className="rdr-recovery__requirement"><ShieldCheck size={17}/> Use at least 12 characters to keep your account secure.</p>
        </div>
        <div className="rdr-recovery__actions">
          <button type="submit" className="rdr-login-v2__submit" disabled={view === 'saving'}>
            {view === 'saving' ? 'Securing your account…' : 'Update password'}
            {view === 'saving' && <LoaderCircle className="rdr-login-v2__spinner" size={18}/>}}
          </button>
          <span className="rdr-recovery__remaining">{invitation && attempts === 1 ? '1 activation attempt remaining' : '\u00a0'}</span>
          <div className="rdr-recovery__error-slot" aria-live="polite">
            {error && <p role="alert" className="rdr-recovery__inline-error"><AlertCircle size={16}/>{error}</p>}
          </div>
        </div>
      </form> : <div className="rdr-recovery__state">
        <RecoveryArtwork state={view} invitation={invitation}/>
        {view === 'checking' && <p className="rdr-recovery__slow" role="status">{slow ? 'Still checking. This may take a moment…' : invitation ? 'Verifying your invitation securely…' : 'Verifying your reset link securely…'}</p>}
        {view === 'success' && <Link className="rdr-login-v2__submit rdr-recovery__state-button" to="/login">Go to login</Link>}
        {view === 'unavailable' && <button className="rdr-login-v2__submit rdr-recovery__state-button" type="button"
          onClick={() => {
            const sequence = ++flowSequence.current;
            if (proofRef.current) void verify(proofRef.current, sequence);
            else {
              const saved = readRecoveryFlow(purpose);
              if (saved?.stage === 'grant') {
                setView('checking');
                void (async () => {
                  try {
                    if (invitation) await accountsApi.setupStatus();
                    else await authenticationApi.resetGrantStatus();
                    if (sequence === flowSequence.current) setView('ready');
                  } catch {
                    if (sequence === flowSequence.current) setView('unavailable');
                  }
                })();
              }
            }
          }}>Try checking again</button>}
        {!invitation && (view === 'missing' || view === 'expired') && <Link className="rdr-login-v2__submit rdr-recovery__state-button" to="/forgot-password">Request new link</Link>}
        {(view === 'missing' || view === 'expired' || view === 'attempts-exhausted') && <Link className="rdr-recovery__back" to="/login">
          <ArrowLeft size={16}/> Return to login
        </Link>}
      </div>}
    </div>
  </LoginLayout>;
};
