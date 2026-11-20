import { useId, useState, type FormEvent } from 'react';
import { ArrowLeft, ArrowRight, Eye, EyeOff, GraduationCap, LockKeyhole, Mail, Monitor, Moon, Sun } from 'lucide-react';
import { Link, useNavigate } from 'react-router-dom';

import { ApiError } from '@Business/api';
import { SchoolMonogram } from '@Components';
import { useWorkspaceTheme, type ThemePreference } from '@Components/hooks/theme';

import { useAuth } from '../core/providers/AuthProvider';

import './loginV2.css';

const emailPattern=/^[^\s@]+@[^\s@]+\.[^\s@]+$/;
const appearanceOptions:{value:ThemePreference;label:string;Icon:typeof Sun}[]=[
  {value:'system',label:'System',Icon:Monitor},
  {value:'light',label:'Light',Icon:Sun},
  {value:'dark',label:'Dark',Icon:Moon},
];

/** Self-contained gradient and vector artwork: no school/stock photographs,
 * external image URLs, or background network requests. */
const GradientArtwork = () => (
  <div className="rdr-login-v2__art" aria-hidden="true">
    <span className="rdr-login-v2__aurora rdr-login-v2__aurora--azure"/>
    <span className="rdr-login-v2__aurora rdr-login-v2__aurora--amber"/>
    <svg className="rdr-login-v2__light-art" viewBox="0 0 640 850" fill="none" preserveAspectRatio="xMidYMid slice" focusable="false">
      <defs>
        <linearGradient id="rdr-v2-neon" x1="54" y1="310" x2="597" y2="516" gradientUnits="userSpaceOnUse">
          <stop stopColor="#4368DB" stopOpacity=".05"/>
          <stop offset=".22" stopColor="#73C5F8"/>
          <stop offset=".43" stopColor="#E7F6FF"/>
          <stop offset=".62" stopColor="#E39BA4"/>
          <stop offset=".78" stopColor="#FF9B4C"/>
          <stop offset="1" stopColor="#FFE9B5" stopOpacity=".04"/>
        </linearGradient>
        <linearGradient id="rdr-v2-ribbon" x1="133" y1="616" x2="497" y2="276" gradientUnits="userSpaceOnUse">
          <stop stopColor="#2974B9" stopOpacity=".1"/>
          <stop offset=".35" stopColor="#65B9F3" stopOpacity=".72"/>
          <stop offset=".67" stopColor="#F6C9A4" stopOpacity=".8"/>
          <stop offset="1" stopColor="#FE7135" stopOpacity=".05"/>
        </linearGradient>
        <linearGradient id="rdr-v2-pages" x1="187" y1="409" x2="486" y2="619" gradientUnits="userSpaceOnUse">
          <stop stopColor="#5CA7F3" stopOpacity=".13"/>
          <stop offset=".49" stopColor="#CFDDFF" stopOpacity=".23"/>
          <stop offset=".92" stopColor="#FFAA65" stopOpacity=".34"/>
        </linearGradient>
        <radialGradient id="rdr-v2-sun" cx="0" cy="0" r="1" gradientTransform="translate(347 411) rotate(90) scale(145 135)" gradientUnits="userSpaceOnUse">
          <stop stopColor="#FFF3C6" stopOpacity=".96"/>
          <stop offset=".14" stopColor="#FFB46A" stopOpacity=".78"/>
          <stop offset=".52" stopColor="#F46D42" stopOpacity=".34"/>
          <stop offset="1" stopColor="#D9555A" stopOpacity="0"/>
        </radialGradient>
        <filter id="rdr-v2-blur" x="-50%" y="-80%" width="200%" height="260%">
          <feGaussianBlur stdDeviation="11"/>
        </filter>
        <filter id="rdr-v2-soft-blur" x="-60%" y="-60%" width="220%" height="220%">
          <feGaussianBlur stdDeviation="30"/>
        </filter>
      </defs>
      <ellipse cx="334" cy="430" rx="178" ry="159" fill="url(#rdr-v2-sun)" filter="url(#rdr-v2-soft-blur)" opacity=".8"/>
      <path d="M161 616C219 515 249 469 335 407C386 369 433 312 445 243" stroke="url(#rdr-v2-ribbon)" strokeWidth="72" opacity=".47" filter="url(#rdr-v2-soft-blur)"/>
      <path d="M154 565C173 472 234 390 341 374C402 364 446 384 476 433" stroke="url(#rdr-v2-ribbon)" strokeWidth="40" opacity=".49" filter="url(#rdr-v2-blur)"/>
      <g className="rdr-login-v2__orbit">
        <ellipse cx="322" cy="413" rx="246" ry="113" transform="rotate(-23 322 413)" stroke="url(#rdr-v2-neon)" strokeWidth="26" filter="url(#rdr-v2-blur)" opacity=".9"/>
        <ellipse cx="322" cy="413" rx="246" ry="113" transform="rotate(-23 322 413)" stroke="url(#rdr-v2-neon)" strokeWidth="3.6" opacity=".92"/>
        <ellipse cx="322" cy="413" rx="230" ry="100" transform="rotate(-24 322 413)" stroke="url(#rdr-v2-neon)" strokeWidth="1.1" opacity=".65"/>
        <path d="M78 463C182 340 354 335 540 415" stroke="#E5F4FF" strokeWidth="1.5" opacity=".4"/>
        <path d="M110 505C234 352 418 334 551 372" stroke="#83CBFF" strokeWidth="1.5" opacity=".42"/>
      </g>
      <g className="rdr-login-v2__book-art">
        <path d="M317 481C273 445 220 439 168 462L163 575C225 551 275 568 318 597Z" fill="url(#rdr-v2-pages)" stroke="#8BD0F7" strokeOpacity=".66" strokeWidth="2.6"/>
        <path d="M322 481C371 444 428 443 478 468L478 575C413 554 364 570 322 597Z" fill="url(#rdr-v2-pages)" stroke="#FFD3A0" strokeOpacity=".75" strokeWidth="2.6"/>
        <path d="M319 481L320 597" stroke="#F2F7FF" strokeWidth="2.5" strokeOpacity=".88"/>
        <path d="M185 482C236 469 271 479 306 503M181 504C230 491 268 503 304 522M182 526C237 515 267 525 302 544" stroke="#BDE5FF" strokeOpacity=".25" strokeWidth="1.5"/>
        <path d="M337 503C379 480 420 476 462 489M337 523C379 498 420 500 460 507M337 545C383 521 424 527 460 530" stroke="#FFE1BE" strokeOpacity=".28" strokeWidth="1.5"/>
        <path d="M319 481C321 453 335 441 344 417" stroke="#F9D4B9" strokeWidth="2" strokeOpacity=".65"/>
      </g>
      <g fill="#E6EFFF">
        <circle cx="142" cy="357" r="2.4" opacity=".75"/>
        <circle cx="525" cy="328" r="1.9" opacity=".72"/>
        <circle cx="456" cy="258" r="1.3" opacity=".7"/>
        <circle cx="191" cy="301" r="1.3" opacity=".7"/>
        <circle cx="483" cy="578" r="1.6" opacity=".5"/>
      </g>
      <path d="M201 271V294M189.5 282.5H212.5M494 273V288M486.5 280.5H501.5" stroke="#F9E8CF" strokeWidth="1.4" strokeLinecap="round" opacity=".66"/>
    </svg>
  </div>
);

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
      <section className="rdr-login-v2__visual" aria-label="RDR School abstract light and learning artwork">
        <GradientArtwork/>
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
