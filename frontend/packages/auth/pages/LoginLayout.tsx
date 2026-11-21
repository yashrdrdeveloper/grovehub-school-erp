import type { ReactNode } from 'react';
import { ArrowLeft, GraduationCap, Monitor, Moon, Sun } from 'lucide-react';
import { Link } from 'react-router-dom';

import { SchoolMonogram } from '@Components';
import { useWorkspaceTheme, type ThemePreference } from '@Components/hooks/theme';

import SoffitGradient from './loginv2/SoffitGradient';

import './loginV2.css';
import './mobileAuth.css';

const appearanceOptions: { value: ThemePreference; label: string; Icon: typeof Sun }[] = [
  { value: 'system', label: 'System', Icon: Monitor },
  { value: 'light', label: 'Light', Icon: Sun },
  { value: 'dark', label: 'Dark', Icon: Moon },
];

/** Shared, card-free shell for sign-in and every recovery step. */
export const LoginLayout = ({ children, title, description, mobileDescription, minimalMobile = false, mobileTitle = 'Sign in', eyebrow = 'YOUR SCHOOL WORKSPACE' }: {
  children: ReactNode;
  title: string;
  description: string;
  mobileDescription?: string;
  minimalMobile?: boolean;
  mobileTitle?: string;
  eyebrow?: string;
}) => {
  const { preference, setPreference } = useWorkspaceTheme();

  return (
    <main className={`rdr-login-v2${minimalMobile ? ' rdr-login-v2--mobile-minimal' : ''}`}>
      <div className="rdr-login-v2__frame">
        <section className="rdr-login-v2__visual" aria-label="RDR School abstract gradient artwork">
          <SoffitGradient/>
          <svg className="rdr-login-v2__mobile-school-art" viewBox="0 0 360 210" fill="none" aria-hidden="true" focusable="false">
            <defs>
              <linearGradient id="rdr-hill" x1="72" y1="180" x2="328" y2="212" gradientUnits="userSpaceOnUse">
                <stop stopColor="#057CA0"/><stop offset="1" stopColor="#22C5B8"/>
              </linearGradient>
              <linearGradient id="rdr-school" x1="220" y1="88" x2="290" y2="168" gradientUnits="userSpaceOnUse">
                <stop stopColor="#D8F8FF"/><stop offset="1" stopColor="#59B9D6"/>
              </linearGradient>
            </defs>
            <path d="M0 174C54 154 79 177 117 161S191 148 229 157 304 143 360 128V210H0Z" fill="#06658A" opacity=".62"/>
            <path d="M0 197C61 168 98 200 157 174S264 157 360 152V210H0Z" fill="url(#rdr-hill)"/>
            <g opacity=".9">
              <circle cx="212" cy="152" r="23" fill="#067D94"/>
              <circle cx="190" cy="160" r="17" fill="#0A8A9C"/>
              <circle cx="331" cy="146" r="20" fill="#18B5B0"/>
            </g>
            <g stroke="#6ED5ED" strokeWidth="2" strokeLinejoin="round">
              <path d="M232 94L272 59 313 94Z" fill="#B7EDF8"/>
              <path d="M237 94V157H308V94Z" fill="url(#rdr-school)"/>
              <path d="M246 100V154M259 100V154M272 100V154M285 100V154M298 100V154" stroke="#E8FDFF" strokeWidth="7"/>
              <path d="M238 157H308M235 162H312M232 167H315" stroke="#B9F1FA" strokeWidth="4"/>
              <path d="M225 116H237V156H217V131Z" fill="#73C4DD" strokeWidth="1.5"/>
              <path d="M219 133H233M219 142H233M219 151H233" stroke="#E1FBFF" strokeWidth="3"/>
              <path d="M308 114H327L337 128V157H308Z" fill="#60B5CF" strokeWidth="1.5"/>
              <path d="M314 132H330M314 142H330M314 152H330" stroke="#E1FBFF" strokeWidth="3"/>
              <path d="M265 157V128Q272 117 280 128V157" fill="#074E75" stroke="#DDFBFF" strokeWidth="3"/>
              <circle cx="272" cy="82" r="6.5" fill="#F2FFFF" stroke="#4BA6C6" strokeWidth="1.5"/>
              <path d="M272 77V82L275 84" stroke="#2A799A" strokeWidth="1.25" strokeLinecap="round"/>
              <path d="M272 60V43M272 43L291 48 272 54" stroke="#DDFBFF" strokeWidth="2.5"/>
            </g>
            <path d="M0 204C77 177 118 210 194 184S308 186 360 171V210H0Z" fill="#20B6B5" opacity=".9"/>
          </svg>
          <div className="rdr-login-v2__mobile-header">
            <div className="rdr-login-v2__mobile-identity" aria-label="RDR School">
              <SchoolMonogram />
              <span className="rdr-login-v2__mobile-school-name">RDR School</span>
            </div>
            <div className="rdr-login-v2__mobile-appearance" role="group" aria-label="Appearance">
              {appearanceOptions.map(({ value, label, Icon }) => (
                <button key={value} type="button" className={preference === value ? 'is-active' : ''}
                  aria-label={`${label} mobile theme`} aria-pressed={preference === value}
                  title={`${label} theme`} onClick={() => setPreference(value)}>
                  <Icon size={16}/>
                </button>
              ))}
            </div>
          </div>
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
        <section className="rdr-login-v2__content" aria-labelledby="rdr-auth-title">
          <div className="rdr-login-v2__toolbar">
            <Link to="/" className="rdr-login-v2__back" aria-label="Back to RDR School homepage">
              <ArrowLeft size={19}/><span>Back to school</span>
            </Link>
            <div className="rdr-login-v2__appearance" role="group" aria-label="Appearance">
              {appearanceOptions.map(({ value, label, Icon }) => (
                <button key={value} type="button" className={preference === value ? 'is-active' : ''}
                  aria-label={`${label} theme`} aria-pressed={preference === value}
                  title={`${label} theme`} onClick={() => setPreference(value)}>
                  <Icon size={16}/>
                </button>
              ))}
            </div>
          </div>
          <div className="rdr-login-v2__main">
            <div className="rdr-login-v2__heading">
              <span className="rdr-login-v2__eyebrow"><GraduationCap size={15}/>{eyebrow}</span>
              <h1 id="rdr-auth-title">{title}</h1>
              <p className="rdr-login-v2__description-desktop">{description}</p>
              <p className="rdr-login-v2__description-mobile">{mobileDescription ?? description}</p>
            </div>
            {minimalMobile && <><h2 className="rdr-login-v2__mobile-signin-title">{mobileTitle}</h2><p className="rdr-login-v2__mobile-description">{mobileDescription ?? description}</p></>}
            {children}
            <div className="rdr-login-v2__closing">
              <span className="rdr-login-v2__closing-rule"/>
              <span>RDR SCHOOL · REWARI</span>
              <span className="rdr-login-v2__closing-rule"/>
            </div>
          </div>
        </section>
      </div>
    </main>
  );
};
