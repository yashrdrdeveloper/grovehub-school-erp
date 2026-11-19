import { useState, type FormEvent } from 'react';
import { Eye, EyeOff, LockKeyhole, LogIn, Mail } from 'lucide-react';
import { Link, useNavigate } from 'react-router-dom';

import { ApiError } from '@Business/api';
import { Button, Input } from '@Components';

import { useAuth } from '../core/providers/AuthProvider';
import { AuthFrame } from './AuthFrame';

const emailPattern = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

export const LoginPage = () => {
  const navigate = useNavigate();
  const { signIn } = useAuth();
  const [identifier, setIdentifier] = useState('');
  const [password, setPassword] = useState('');
  const [visible, setVisible] = useState(false);
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState('');
  const [fieldErrors, setFieldErrors] = useState<{ identifier?: string; password?: string }>({});

  const submit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setError('');
    const next: typeof fieldErrors = {};
    const email = identifier.trim();
    if (!email) next.identifier = 'Enter your email address.';
    else if (!emailPattern.test(email)) next.identifier = 'Enter a valid email address.';
    if (!password) next.password = 'Enter your password.';
    setFieldErrors(next);
    if (Object.keys(next).length) return;

    setBusy(true);
    try {
      await signIn(email, password);
      navigate('/dashboard', { replace: true });
    } catch (failure) {
      setError(failure instanceof ApiError ? failure.message : 'Unable to sign in. Please try again.');
    } finally {
      setBusy(false);
    }
  };

  return (
    <AuthFrame title="Sign in">
      <div className="auth-panel__icon"><LogIn size={23}/></div>
      <h1>Sign in</h1>
      <p className="muted"><span className="auth-login__desktop-copy">Enter your school account credentials.</span><span className="auth-login__mobile-copy">Welcome back. Sign in to your school workspace.</span></p>
      <form noValidate onSubmit={(event) => void submit(event)} className="auth-form">
        <div className="auth-login__input-wrap">
        <Input
          id="identifier"
          type="email"
          inputMode="email"
          label="Email"
          placeholder="Enter your email"
          value={identifier}
          autoComplete="username"
          error={fieldErrors.identifier}
          onChange={(event) => {
            setIdentifier(event.target.value);
            if (fieldErrors.identifier) setFieldErrors(current => ({ ...current, identifier: undefined }));
            if (error) setError('');
          }}
          required
        />
          <Mail className="auth-login__input-icon" size={18} aria-hidden="true"/>
        </div>
        <div className="password-field">
          <Input
            id="password"
            label="Password"
            placeholder="Enter your password"
            type={visible ? 'text' : 'password'}
            value={password}
            autoComplete="current-password"
            error={fieldErrors.password}
            onChange={(event) => {
              setPassword(event.target.value);
              if (fieldErrors.password) setFieldErrors(current => ({ ...current, password: undefined }));
              if (error) setError('');
            }}
            required
          />
          <LockKeyhole className="auth-login__input-icon" size={18} aria-hidden="true"/>
          <button type="button" aria-label={visible ? 'Hide password' : 'Show password'}
            onClick={() => setVisible(current => !current)} className="password-toggle">
            {visible ? <EyeOff size={18} /> : <Eye size={18} />}
          </button>
        </div>
        {error ? <p className="form-error" role="alert">{error}</p> : null}
        <div className="form-row form-row--end"><Link to="/forgot-password">Forgot password?</Link></div>
        <Button type="submit" loading={busy} className="button--full">Continue</Button>
      </form>
      <p className="auth-security"><LockKeyhole size={14} /> Protected school workspace</p>
      <p className="auth-login__contact">Need an account? <Link to="/contact">Contact school</Link></p>
    </AuthFrame>
  );
};
