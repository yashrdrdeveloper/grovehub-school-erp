import { useId, useState, type FormEvent } from 'react';
import { ArrowLeft, ArrowRight, Eye, EyeOff, GraduationCap, LockKeyhole, Mail, Monitor, Moon, Sun } from 'lucide-react';
import { Link, useNavigate } from 'react-router-dom';

import { ApiError } from '@Business/api';
import { SchoolMonogram } from '@Components';
import { useWorkspaceTheme, type ThemePreference } from '@Components/hooks/theme';

import { useAuth } from '../core/providers/AuthProvider';

import './loginV2.css';

const emailPattern=/^[^\s@]+@[^\s@]+\.[^\s@]+$/;
const schoolPhoto='/school-images/slide5.jpg';
const originalSchoolPhoto='https://rdrschool.com/images/slide5.jpg';

const appearanceOptions:{value:ThemePreference;label:string;Icon:typeof Sun}[]=[
  {value:'system',label:'System',Icon:Monitor},
  {value:'light',label:'Light',Icon:Sun},
  {value:'dark',label:'Dark',Icon:Moon},
];

export const LoginV2Page = () => {
  const fieldId=useId();
  const navigate=useNavigate();
  const {signIn}=useAuth();
  const {preference,setPreference}=useWorkspaceTheme();

  const [email,setEmail]=useState('');
  const [password,setPassword]=useState('');
  const [visible,setVisible]=useState(false);
  const [busy,setBusy]=useState(false);
  const [error,setError]=useState('');
  const [fieldErrors,setFieldErrors]=useState<{email?:string;password?:string}>({});
  const [photoSrc,setPhotoSrc]=useState(schoolPhoto);
  const [photoUnavailable,setPhotoUnavailable]=useState(false);

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
      setError(failure instanceof ApiError?failure.message:'Unable to sign in. Please try again.');
    }finally{
      setBusy(false);
    }
  };

  return <main className="rdr-login-v2">
    <div className="rdr-login-v2__frame">
      <section className="rdr-login-v2__visual" aria-label="RDR School student life">
        {!photoUnavailable&&<img className="rdr-login-v2__photo" src={photoSrc}
          alt="RDR School students together in school uniforms"
          loading="eager" decoding="async"
          onError={()=>{
            if(photoSrc===schoolPhoto)setPhotoSrc(originalSchoolPhoto);
            else setPhotoUnavailable(true);
          }}/>}
        <div className="rdr-login-v2__photo-tint" aria-hidden="true"/>
        <div className="rdr-login-v2__visual-brand">
          <SchoolMonogram/>
          <span><strong>RDR School</strong><small>Rao Deena Ram Vidya Vihar</small></span>
        </div>
        <div className="rdr-login-v2__visual-copy">
          <span className="rdr-login-v2__visual-kicker"><span/> A PLACE TO BELONG</span>
          <h2>Growing minds.<br/><em>Brighter futures.</em></h2>
          <p>Learning, character and community — together at RDR School.</p>
          <div className="rdr-login-v2__visual-footer"><span>HALUHERA</span><i aria-hidden="true"/><span>REWARI, HARYANA</span></div>
        </div>
        <svg className="rdr-login-v2__doodle" aria-hidden="true" viewBox="0 0 190 110" fill="none">
          <path d="M5 83c37-40 51 18 91-27S159 29 183 17" stroke="currentColor" strokeWidth="1.4" strokeDasharray="4 7" strokeLinecap="round"/>
          <path d="m158 18 25-1-13 21-4-11-8-9Zm8 9 17-10" stroke="currentColor" strokeWidth="1.4" strokeLinejoin="round"/>
          <path d="M24 19v12M18 25h12M44 6v9M40 10.5h9" stroke="currentColor" strokeLinecap="round" strokeWidth="1.3"/>
        </svg>
      </section>

      <section className="rdr-login-v2__content" aria-labelledby="rdr-login-v2-title">
        <div className="rdr-login-v2__toolbar">
          <Link to="/" className="rdr-login-v2__back" aria-label="Back to RDR School homepage">
            <ArrowLeft size={19}/><span>Back to school</span>
          </Link>
          <div className="rdr-login-v2__appearance" role="group" aria-label="Appearance">
            {appearanceOptions.map(({value,label,Icon})=>
              <button key={value} type="button" className={preference===value?'is-active':''}
                aria-label={`${label} theme`} aria-pressed={preference===value}
                title={`${label} theme`} onClick={()=>setPreference(value)}>
                <Icon size={16}/>
              </button>
            )}
          </div>
        </div>

        <div className="rdr-login-v2__main">
          <div className="rdr-login-v2__heading">
            <span className="rdr-login-v2__eyebrow"><GraduationCap size={15}/> YOUR SCHOOL WORKSPACE</span>
            <h1 id="rdr-login-v2-title">Welcome back.</h1>
            <p>Sign in to continue to your RDR School workspace.</p>
          </div>

          <form className="rdr-login-v2__form" onSubmit={event=>void submit(event)} noValidate>
            <div className="rdr-login-v2__field">
              <label htmlFor={`${fieldId}-email`}>Email address</label>
              <div className="rdr-login-v2__input-shell">
                <Mail size={18} aria-hidden="true"/>
                <input id={`${fieldId}-email`} type="email" inputMode="email"
                  autoComplete="username" placeholder="you@school.edu" value={email}
                  aria-invalid={Boolean(fieldErrors.email)}
                  aria-describedby={fieldErrors.email?`${fieldId}-email-error`:undefined}
                  onChange={event=>{
                    setEmail(event.target.value);
                    if(fieldErrors.email)setFieldErrors(current=>({...current,email:undefined}));
                    if(error)setError('');
                  }} required/>
              </div>
              {fieldErrors.email&&<small className="rdr-login-v2__field-error" id={`${fieldId}-email-error`}>{fieldErrors.email}</small>}
            </div>

            <div className="rdr-login-v2__field">
              <div className="rdr-login-v2__password-label">
                <label htmlFor={`${fieldId}-password`}>Password</label>
                <Link to="/forgot-password">Forgot password?</Link>
              </div>
              <div className="rdr-login-v2__input-shell">
                <LockKeyhole size={18} aria-hidden="true"/>
                <input id={`${fieldId}-password`} type={visible?'text':'password'}
                  autoComplete="current-password" placeholder="Enter your password"
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
              {fieldErrors.password&&<small className="rdr-login-v2__field-error" id={`${fieldId}-password-error`}>{fieldErrors.password}</small>}
            </div>

            {error&&<p className="rdr-login-v2__error" role="alert">{error}</p>}
            <button type="submit" className="rdr-login-v2__submit" disabled={busy}>
              <span>{busy?'Signing in…':'Sign in to workspace'}</span><ArrowRight size={19}/>
            </button>
          </form>
          <div className="rdr-login-v2__closing"><span className="rdr-login-v2__closing-rule"/><span>RDR SCHOOL · REWARI</span><span className="rdr-login-v2__closing-rule"/></div>
        </div>
      </section>
    </div>
  </main>;
};
