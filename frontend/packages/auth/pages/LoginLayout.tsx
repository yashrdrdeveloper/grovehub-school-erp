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
          <div className="rdr-login-v2__mobile-header">
            <div className="rdr-login-v2__mobile-identity" aria-label="RDR School">
              <SchoolMonogram />
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
            {minimalMobile && <div className="rdr-login-v2__mobile-signin-title">{mobileTitle}</div>}
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
