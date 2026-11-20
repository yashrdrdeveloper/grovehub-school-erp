import type { ReactNode } from 'react';
import { ArrowLeft, BookOpen, CheckCircle2, GraduationCap, Monitor, Moon, ShieldCheck, Sparkles, Sun } from 'lucide-react';
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
export const LoginLayout = ({ children, title, description, eyebrow = 'YOUR SCHOOL WORKSPACE' }: {
  children: ReactNode;
  title: string;
  description: string;
  eyebrow?: string;
}) => {
  const { preference, setPreference } = useWorkspaceTheme();

  return (
    <main className="rdr-login-v2">
      <div className="rdr-login-v2__frame">
        <section className="rdr-login-v2__visual" aria-label="RDR School abstract gradient artwork">
          <SoffitGradient/>
          <div className="rdr-login-v2__mobile-header">
            <div className="rdr-login-v2__mobile-identity">
              <SchoolMonogram />
              <span><strong>RDR School</strong><small>Learn · Grow · Belong</small></span>
            </div>
            <span className="rdr-login-v2__mobile-security" aria-label="Secure authentication"><ShieldCheck size={18}/></span>
          </div>
          <div className="rdr-login-v2__mobile-illustration" aria-hidden="true">
            <div className="rdr-login-v2__mobile-glow"/>
            <div className="rdr-login-v2__mobile-art-card rdr-login-v2__mobile-art-card--back"><Sparkles size={21}/><span className="rdr-login-v2__mobile-art-lines"><i/><i/><i/></span></div>
            <div className="rdr-login-v2__mobile-art-card rdr-login-v2__mobile-art-card--front"><BookOpen size={30}/><span><i/><i/></span></div>
            <div className="rdr-login-v2__mobile-art-check"><CheckCircle2 size={23}/></div>
            <span className="rdr-login-v2__mobile-art-dot rdr-login-v2__mobile-art-dot--one"/>
            <span className="rdr-login-v2__mobile-art-dot rdr-login-v2__mobile-art-dot--two"/>
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
              <p>{description}</p>
            </div>
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
