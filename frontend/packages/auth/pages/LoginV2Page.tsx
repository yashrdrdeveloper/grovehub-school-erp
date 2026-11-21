import { useId, useState, type FormEvent } from 'react';
import { AlertCircle, ArrowRight, Eye, EyeOff, LoaderCircle, LockKeyhole, Mail } from 'lucide-react';
import { Link, useNavigate } from 'react-router-dom';

import { useAuth } from '../core/providers/AuthProvider';

import { LoginLayout } from './LoginLayout';

const emailPattern=/^[^\s@]+@[^\s@]+\.[^\s@]+$/;
export const LoginV2Page = () => {
  const fieldId=useId();
  const navigate=useNavigate();
  const {signIn}=useAuth();

  const [email,setEmail]=useState('');
  const [password,setPassword]=useState('');
  const [visible,setVisible]=useState(false);
  const [busy,setBusy]=useState(false);
  const [error,setError]=useState('');
  const [fieldErrors,setFieldErrors]=useState<{email?:string;password?:string}>({});

  const submit=async(event:FormEvent<HTMLFormElement>)=>{
    event.preventDefault();
    if(busy)return;
    const normalizedEmail=email.trim();
    const nextErrors:typeof fieldErrors={};
    if(!normalizedEmail)nextErrors.email='Please enter your email address.';
    else if(!emailPattern.test(normalizedEmail))nextErrors.email='Enter a valid email address.';
    if(!password)nextErrors.password='Please enter your password.';
    setFieldErrors(nextErrors);
    setError('');
    if(Object.keys(nextErrors).length)return;
    setBusy(true);
    try{
      await signIn(normalizedEmail,password);
      navigate('/dashboard',{replace:true});
    }catch(failure){
      const message=failure instanceof Error?failure.message:'';
      const invalidCredentials=/invalid creden|incorrect|unauthoriz|invalid login|wrong password|\b401\b/i.test(message);
      setError(invalidCredentials
        ? 'Incorrect email or password. Please try again.'
        : 'Unable to sign in right now. Please try again.');
    }finally{
      setBusy(false);
    }
  };

  return <LoginLayout title="Welcome back." description="Sign in to continue to your RDR School workspace." mobileDescription="Sign in to your workspace." minimalMobile>
          <form className="rdr-login-v2__form rdr-login-v2__form--signin" onSubmit={event=>void submit(event)} noValidate>
            <div className="rdr-login-v2__field">
              <label htmlFor={`${fieldId}-email`}>Email address</label>
              <div className="rdr-login-v2__input-shell">
                <Mail size={18} aria-hidden="true"/>
                <input id={`${fieldId}-email`} type="email" inputMode="email"
                  autoComplete="username" placeholder="Email" value={email}
                  aria-invalid={Boolean(fieldErrors.email)}
                  aria-describedby={fieldErrors.email?`${fieldId}-email-error`:undefined}
                  onChange={event=>{
                    setEmail(event.target.value);
                    if(fieldErrors.email)setFieldErrors(current=>({...current,email:undefined}));
                    if(error)setError('');
                  }} required/>
              </div>
              <small className={`rdr-login-v2__field-error rdr-login-v2__field-error--reserved${fieldErrors.email ? "" : " is-empty"}`} id={`${fieldId}-email-error`}>{fieldErrors.email || '\u00a0'}</small>
            </div>

            <div className="rdr-login-v2__field">
              <div className="rdr-login-v2__password-label">
                <label htmlFor={`${fieldId}-password`}>Password</label>
                <Link className="rdr-login-v2__forgot-desktop" to="/forgot-password">Forgot password?</Link>
              </div>
              <div className="rdr-login-v2__input-shell">
                <LockKeyhole size={18} aria-hidden="true"/>
                <input id={`${fieldId}-password`} type={visible?'text':'password'}
                  autoComplete="current-password" placeholder="Password"
                  value={password} aria-invalid={Boolean(fieldErrors.password)}
                  aria-describedby={fieldErrors.password?`${fieldId}-password-error`:undefined}
                  onChange={event=>{
                    setPassword(event.target.value);
                    if(fieldErrors.password)setFieldErrors(current=>({...current,password:undefined}));
                    if(error)setError('');
                  }} required/>
                <button type="button" className="rdr-login-v2__visibility"
                  onClick={()=>setVisible(current=>!current)}
                  aria-label={visible?'Hide password':'Show password'}>
                  {visible?<EyeOff size={19}/>:<Eye size={19}/>}
                </button>
              </div>
              <small className={`rdr-login-v2__field-error rdr-login-v2__field-error--reserved${fieldErrors.password ? "" : " is-empty"}`} id={`${fieldId}-password-error`}>{fieldErrors.password || '\u00a0'}</small>
            </div>

            <div className="rdr-login-v2__feedback-slot rdr-login-v2__feedback-slot--signin" aria-live="polite">
              {error && <p className="rdr-login-v2__auth-error" role="alert">
                <AlertCircle size={16} aria-hidden="true"/>
                <span>{error}</span>
              </p>}
            </div>
            <Link className="rdr-login-v2__forgot-mobile" to="/forgot-password">Forgot password?</Link>
            <button type="submit" className="rdr-login-v2__submit" disabled={busy}>
              <span>{busy?'Signing in…':'Sign in'}</span>
              {busy?<LoaderCircle className="rdr-login-v2__spinner" size={19} aria-hidden="true"/>:<ArrowRight size={19} aria-hidden="true"/>}
            </button>
          </form>
  </LoginLayout>;
};
