import { useState, type FormEvent } from 'react';
import SoffitGradient from './SoffitGradient';
import './LoginV2.css';

export type LoginV2Credentials = { email: string; password: string; remember: boolean };
export type LoginV2Props = {
  onSignIn: (credentials: LoginV2Credentials) => Promise<void> | void;
  onForgotPassword?: () => void;
  schoolName?: string;
  schoolSubtitle?: string;
};

export default function LoginV2({ onSignIn, onForgotPassword, schoolName = 'GrooveHub School', schoolSubtitle = 'Your campus, beautifully connected.' }: LoginV2Props) {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [remember, setRemember] = useState(true);
  const [visible, setVisible] = useState(false);
  const [theme, setTheme] = useState<'light' | 'dark'>(() => typeof window !== 'undefined' && window.matchMedia?.('(prefers-color-scheme: dark)').matches ? 'dark' : 'light');
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState('');
  async function submit(e: FormEvent<HTMLFormElement>) {
    e.preventDefault();
    if (busy) return;
    setBusy(true);
    setError('');
    try { await onSignIn({ email: email.trim(), password, remember }); }
    catch (err) { setError(err instanceof Error ? err.message : 'Unable to sign in. Please try again.'); }
    finally { setBusy(false); }
  }
  return (
    <main className="lv2" data-theme={theme}>
      <div className="lv2-layout">
        <section className="lv2-hero" aria-label="Welcome to GrooveHub">
          <SoffitGradient />
          <div className="lv2-hero-inner">
            <div className="lv2-brand"><span className="lv2-symbol">✳</span> groove<span>hub</span></div>
            <div className="lv2-story">
              <div className="lv2-kicker">THE MODERN SCHOOL EXPERIENCE</div>
              <h1>One place for <em>every school</em> moment.</h1>
              <p>Where learning, people and possibilities come together — thoughtfully designed for your school community.</p>
              <div className="lv2-pills"><span>✦ Connected learning</span><span>◉ Made for everyone</span></div>
              <div className="lv2-float"><span>✧</span><div><strong>Every day, in sync</strong><small>One connected campus</small></div></div>
            </div>
            <div className="lv2-hero-foot"><span>Built for growing minds.</span><span>© GrooveHub</span></div>
          </div>
        </section>
        <section className="lv2-panel" aria-labelledby="lv2-title">
          <header className="lv2-panel-head">
            <div className="lv2-mobile-brand"><span className="lv2-symbol">✳</span> groove<strong>hub</strong></div>
            <button type="button" className="lv2-mode" onClick={() => setTheme(theme === 'light' ? 'dark' : 'light')} aria-label={`Switch to ${theme === 'light' ? 'dark' : 'light'} mode`}>{theme === 'light' ? '☾ Dark' : '☀ Light'} mode</button>
          </header>
          <div className="lv2-form-wrap">
            <div className="lv2-mobile-banner" aria-hidden="true"><SoffitGradient /><span>YOUR SCHOOL, YOUR SPACE ✳</span></div>
            <div className="lv2-school"><div className="lv2-school-mark">⌂</div><div><strong>{schoolName}</strong><small>{schoolSubtitle}</small></div></div>
            <div className="lv2-intro"><span>WELCOME BACK ✳</span><h2 id="lv2-title">Good to see you again.</h2><p>Sign in to pick up right where you left off.</p></div>
            <form className="lv2-form" onSubmit={submit}>
              <label htmlFor="lv2-email">Email address</label>
              <input id="lv2-email" type="email" inputMode="email" autoComplete="username" placeholder="name@school.edu" value={email} onChange={e => setEmail(e.target.value)} required />
              <div className="lv2-row"><label htmlFor="lv2-password">Password</label><button type="button" className="lv2-link" onClick={onForgotPassword} disabled={!onForgotPassword}>Forgot password?</button></div>
              <div className="lv2-password"><input id="lv2-password" type={visible ? 'text' : 'password'} autoComplete="current-password" placeholder="Enter your password" value={password} onChange={e => setPassword(e.target.value)} required /><button type="button" aria-label={visible ? 'Hide password' : 'Show password'} onClick={() => setVisible(!visible)}>{visible ? 'Hide' : 'Show'}</button></div>
              <label className="lv2-remember"><input type="checkbox" checked={remember} onChange={e => setRemember(e.target.checked)} />Keep me signed in</label>
              {error && <p className="lv2-error" role="alert">{error}</p>}
              <button className="lv2-submit" type="submit" disabled={busy}>{busy ? 'Signing in…' : 'Continue to your school'} <span aria-hidden="true">↗</span></button>
            </form>
            <div className="lv2-secure">◇ A private, secure space for your school</div>
          </div>
          <footer className="lv2-bottom">Made for the moments that matter. <span>© GrooveHub</span></footer>
        </section>
      </div>
    </main>
  );
}
