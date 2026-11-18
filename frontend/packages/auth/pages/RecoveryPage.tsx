import { useEffect, useState, type FormEvent } from 'react';
import { CheckCircle2, Eye, EyeOff } from 'lucide-react';
import { Link } from 'react-router-dom';

import { authenticationApi } from '@Business/auth';
import { accountsApi } from '@Business/accounts';
import { Button, Input } from '@Components';

import { AuthFrame } from './AuthFrame';

type Purpose = 'recovery' | 'invitation';

export const RecoveryPage = ({ purpose }: { purpose: Purpose }) => {
  const [verified, setVerified] = useState(false);
  const [password, setPassword] = useState('');
  const [confirmation, setConfirmation] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [username, setUsername] = useState('');
  const [complete, setComplete] = useState(false);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    // URL fragments are never sent in an HTTP request or referrer.
    const params = new URLSearchParams(window.location.hash.slice(1));
    const proof = params.get('proof');
    window.history.replaceState(null, '', window.location.pathname);
    if (!proof) {
      setError('The link is missing its one-time proof. Request a new email.');
      return;
    }
    const action = purpose === 'recovery'
      ? authenticationApi.exchangeRecovery(proof)
      : authenticationApi.exchangeInvitation(proof);
    void action.then(
      () => setVerified(true),
      () => setError('This link is expired or already used. Request a fresh one.'),
    );
  }, [purpose]);

  const submit = async (event: FormEvent) => {
    event.preventDefault();
    setError('');
    if (password !== confirmation) {
      setError('Passwords do not match.');
      return;
    }
    setBusy(true);
    try {
      if (purpose === 'recovery') {
        await authenticationApi.resetPassword(password);
      } else {
        await accountsApi.setup(username.trim(), password);
      }
      setComplete(true);
    } catch (failure) {
      setError(failure instanceof Error ? failure.message : 'Unable to complete setup.');
    } finally {
      setBusy(false);
    }
  };

  const settingUp = purpose === 'invitation';

  return (
    <AuthFrame title={settingUp ? "Set password" : "Reset password"}>
                <h1>{settingUp ? 'Set up your account' : 'Choose a new password'}</h1>
        <p className="muted">
          {settingUp ? 'Complete the invitation to activate your school account.' : 'Your new password will replace the previous one.'}
        </p>
        {complete ? (
          <div className="notice"><CheckCircle2 size={20} /> All set. <Link to="/login">Go to sign in</Link></div>
        ) : verified ? (
          <form className="auth-form" onSubmit={(event) => void submit(event)}>
            {settingUp ? <Input id="username" label="Username" placeholder="Choose a username" autoComplete="username" value={username} onChange={(event) => setUsername(event.target.value)} required /> : null}
            <div className="password-field">
              <Input id="newPassword" label="New password" placeholder="Create a secure password" type={showPassword ? 'text' : 'password'} value={password} minLength={12} autoComplete="new-password" onChange={(event) => setPassword(event.target.value)} required hint="At least 12 characters." />
              <button type="button" className="password-toggle" onClick={() => setShowPassword(value => !value)} aria-label={showPassword ? 'Hide passwords' : 'Show passwords'}>{showPassword ? <EyeOff size={18}/> : <Eye size={18}/>}</button>
            </div>
            <Input id="confirmPassword" label="Confirm new password" placeholder="Re-enter your new password" type={showPassword ? 'text' : 'password'} value={confirmation} minLength={12} autoComplete="new-password" onChange={event => setConfirmation(event.target.value)} required />
            {error ? <p className="form-error" role="alert">{error}</p> : null}
            <Button type="submit" loading={busy} className="button--full">{settingUp ? 'Activate account' : 'Reset password'}</Button>
          </form>
        ) : <div className="notice" role="status">{error || 'Verifying your secure link…'}</div>}
    </AuthFrame>
  );
};
