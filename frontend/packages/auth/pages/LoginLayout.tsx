import { useEffect, type ReactNode } from 'react';
import { ArrowLeft, LockKeyhole, Monitor, Moon, Sun } from 'lucide-react';
import { Link } from 'react-router-dom';

import { SchoolMonogram } from '@Components';
import { useWorkspaceTheme, type ThemePreference } from '@Components/hooks/theme';

import { AuthCampusIllustration } from './AuthCampusIllustration';


// One stylesheet owns the entire auth layout. Legacy variants are intentionally not loaded.
import './optionOneAuth.css';

const appearanceOptions: { value: ThemePreference; label: string; Icon: typeof Sun }[] = [
  { value: 'system', label: 'System', Icon: Monitor },
  { value: 'light', label: 'Light', Icon: Sun },
  { value: 'dark', label: 'Dark', Icon: Moon },
];

/** Shared, card-free shell for sign-in and every recovery step. */
export const LoginLayout = ({ children, title, description, mobileDescription, minimalMobile = false, eyebrow = 'YOUR SCHOOL WORKSPACE' }: {
  children: ReactNode;
  title: string;
  description: string;
  mobileDescription?: string;
  minimalMobile?: boolean;
  mobileTitle?: string;
  eyebrow?: string;
}) => {
  const { preference, setPreference } = useWorkspaceTheme();

  // On phones use the operating system's light/dark preference.
  // The existing desktop theme selector remains available.
  useEffect(() => {
    if (typeof window.matchMedia !== 'function') return;
    const mobile = window.matchMedia('(max-width: 720px)');
    const syncMobileAppearance = () => {
      if (mobile.matches && preference !== 'system') setPreference('system');
    };
    syncMobileAppearance();
    mobile.addEventListener('change', syncMobileAppearance);
    return () => mobile.removeEventListener('change', syncMobileAppearance);
  }, [preference, setPreference]);

  return (
    <main className={`rdr-login-v2${minimalMobile ? ' rdr-login-v2--mobile-minimal' : ''}`}>
      <div className="rdr-login-v2__frame">
        <section className="rdr-login-v2__visual" aria-label="RDR School branded gradient background">
          <div className="rdr-login-v2__mobile-header">
            <div className="rdr-login-v2__mobile-identity" aria-label="RDR School">
              <SchoolMonogram/>
              <span className="rdr-login-v2__mobile-school-name">RDR School<small>Your school workspace</small></span>
            </div>
          </div>
          <div className="rdr-login-v2__visual-brand">
            <SchoolMonogram/>
            <span><strong>RDR School</strong><small>Rao Deena Ram Vidya Vihar</small></span>
          </div>
          <AuthCampusIllustration/>
          <div className="rdr-login-v2__visual-copy">
            <span className="rdr-login-v2__visual-kicker"><span/> THE RDR SCHOOL WORKSPACE</span>
            <h2>For every school day.<br/><em>And everyone behind it.</em></h2>
            <p>A connected space for the people who make learning happen.</p>
            <div className="rdr-login-v2__visual-footer"><span>LEARN</span><i aria-hidden="true"/><span>ORGANISE</span><i aria-hidden="true"/><span>GROW</span></div>
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
              <span className="rdr-login-v2__eyebrow">{eyebrow}</span>
              <h1 id="rdr-auth-title">{title}</h1>
              <p className="rdr-login-v2__description-desktop">{description}</p>
              <p className="rdr-login-v2__description-mobile">{mobileDescription ?? description}</p>
            </div>
            {children}
            <div className="rdr-login-v2__closing" aria-label="Secure RDR School workspace">
              <LockKeyhole size={13} strokeWidth={1.8} aria-hidden="true"/>
              <span>A secure workspace for RDR School</span>
            </div>
            <div className="rdr-o1-auth-footer"><span>© RDR School</span><span>Haluhera · Rewari</span></div>
          </div>
        </section>
      </div>
    </main>
  );
};
