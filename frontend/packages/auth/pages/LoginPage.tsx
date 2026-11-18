import { useState, type FormEvent } from 'react';
import { Eye, EyeOff, LockKeyhole } from 'lucide-react';
import { Link, useNavigate } from 'react-router-dom';

import { ApiError } from '@Business/api';
import { Button, Input } from '@Components';

import { useAuth } from '../core/providers/AuthProvider';

export const LoginPage = () => {
  const navigate = useNavigate();
  const { signIn } = useAuth();
  const [identifier, setIdentifier] = useState('');
  const [password, setPassword] = useState('');
  const [visible, setVisible] = useState(false);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState('');

  const submit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError('');
    setBusy(true);
    try {
      await signIn(identifier.trim(), password);
      navigate('/dashboard', { replace: true });
    } catch (failure) {
      setError(failure instanceof ApiError ? failure.message : 'Unable to sign in. Please try again.');
    } finally {
      setBusy(false);
    }
  };

  return (
    <main className="auth-page">
      <div className="auth-brand">
        <div className="brand-mark">G</div>
        <strong>groovehub<span className="brand-period">.</span></strong>
      </div>
      <section className="auth-panel" aria-label="Sign in">
        <div className="auth-panel__school"><span className="school-emblem">✦</span> RDR School</div>
        <h1>Sign in</h1>
        <p className="muted">Enter your school account credentials.</p>
        <form onSubmit={(event) => void submit(event)} className="auth-form">
          <Input
            id="identifier"
            label="Email"
            value={identifier}
            autoComplete="username"
            onChange={(event) => setIdentifier(event.target.value)}
            required
            autoFocus
          />
          <div className="password-field">
            <Input
              id="password"
              label="Password"
              type={visible ? 'text' : 'password'}
              value={password}
              autoComplete="current-password"
              onChange={(event) => setPassword(event.target.value)}
              required
            />
            <button
              type="button"
              aria-label={visible ? 'Hide password' : 'Show password'}
              onClick={() => setVisible(!visible)}
              className="password-toggle"
            >
              {visible ? <EyeOff size={18} /> : <Eye size={18} />}
            </button>
          </div>
          {error ? <p className="form-error" role="alert">{error}</p> : null}
          <div className="form-row form-row--end"><Link to="/forgot-password">Forgot password?</Link></div>
          <Button type="submit" loading={busy} className="button--full">Continue</Button>
        </form>
        <p className="auth-security"><LockKeyhole size={14} /> Protected school workspace</p>
      </section>
      <footer className="auth-footer">
        <span>© {new Date().getFullYear()} GrooveHub</span>
        
      </footer>
    </main>
  );
};
