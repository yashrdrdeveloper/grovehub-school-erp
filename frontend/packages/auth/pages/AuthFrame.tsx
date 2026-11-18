import type { ReactNode } from 'react';
import {
  ArrowRight, BookOpen, CalendarDays, GraduationCap, ListTodo, Monitor,
  Moon, Sparkles, Sun, UsersRound, WalletCards,
} from 'lucide-react';

import { SchoolMonogram } from '@Components';
import { useWorkspaceTheme, type ThemePreference } from '@Shell/theme';

const themeChoices: { value: ThemePreference; label: string; icon: typeof Sun }[] = [
  { value: 'system', label: 'System', icon: Monitor },
  { value: 'light', label: 'Light', icon: Sun },
  { value: 'dark', label: 'Dark', icon: Moon },
];

/** All artwork stays inside its own layout lane, never above the sign-in form. */
const SchoolIllustrationLeft = () => (
  <div className="auth-art-lane auth-art-lane--left" aria-hidden="true">
    <div className="auth-art-doodle auth-art-doodle--calendar"><CalendarDays size={26} strokeWidth={1.35}/></div>
    <svg className="auth-art-squiggle" viewBox="0 0 140 50" fill="none">
      <path d="M2 30c16 30 28-30 44-11 12 16 0 31 19 23 22-10 22-43 42-24 12 13 18 19 31 12" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round"/>
    </svg>
    <blockquote className="auth-art-quote">Everything your school needs, in one calm workspace.</blockquote>
    <div className="auth-art-card auth-art-card--crm">
      <div className="auth-art-card__title"><UsersRound size={18}/><strong>Lead CRM</strong></div>
      <div className="auth-art-card__rule"/>
      <div className="auth-art-card__journey"><span>Enquiry</span><ArrowRight size={13}/><span>Admission</span></div>
      <div className="auth-art-card__contact"><i/><b/><span>✓</span></div>
      <div className="auth-art-card__contact"><i/><b/><span>✓</span></div>
    </div>
    <div className="auth-art-card auth-art-card--small">
      <ListTodo size={18}/>
      <span><strong>Tasks</strong><small>Daily school work</small></span>
    </div>
    <div className="auth-art-lane__footer">LEARN <i>·</i> ORGANIZE <i>·</i> GROW</div>
  </div>
);

const SchoolIllustrationRight = () => (
  <div className="auth-art-lane auth-art-lane--right" aria-hidden="true">
    <div className="auth-art-doodle auth-art-doodle--book"><BookOpen size={27} strokeWidth={1.35}/></div>
    <svg className="auth-art-squiggle" viewBox="0 0 140 50" fill="none">
      <path d="M2 30c16 30 28-30 44-11 12 16 0 31 19 23 22-10 22-43 42-24 12 13 18 19 31 12" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round"/>
    </svg>
    <div className="auth-art-card auth-art-card--student">
      <div className="auth-art-card__title"><GraduationCap size={18}/><strong>Student Lifecycle 360°</strong></div>
      <small>Unified student view · UCV</small>
      <div className="auth-art-card__rule"/>
      <div className="auth-art-card__student">
        <span className="auth-art-card__avatar"/>
        <span className="auth-art-card__lines"><i/><i/></span>
        <span className="auth-art-card__ring"/>
      </div>
      <div className="auth-art-card__tags"><span>Academics</span><span>Growth</span></div>
    </div>
    <div className="auth-art-card auth-art-card--small">
      <WalletCards size={18}/>
      <span><strong>Finance &amp; Fees</strong><small>School finances</small></span>
    </div>
    <div className="auth-art-lane__footer">RDR SCHOOL <i>/</i> GROOVEHUB</div>
  </div>
);

const MobileIllustration = () => (
  <div className="auth-mobile-art" aria-hidden="true">
    <div className="auth-mobile-art__line"/>
    <div className="auth-mobile-art__story"><Sparkles size={15}/><span>From first enquiry to every milestone.</span></div>
    <div className="auth-mobile-art__motto">LEARN <i>·</i> ORGANIZE <i>·</i> GROW</div>
  </div>
);

export const AuthFrame = ({ children, title }: { children: ReactNode; title: string }) => {
  const { preference, setPreference } = useWorkspaceTheme();

  return (
    <main className="auth-page auth-page--consistent">
      <header className="auth-topbar">
        <div className="auth-brand auth-brand--school" aria-label="RDR School">
          <SchoolMonogram />
          <strong>RDR School</strong>
        </div>
        <div className="auth-theme-switcher" role="group" aria-label="Appearance">
          {themeChoices.map(({ value, label, icon: Icon }) => (
            <button key={value} type="button" title={label} aria-label={`${label} theme`}
              aria-pressed={preference === value}
              className={preference === value ? 'auth-theme-switcher__choice is-active' : 'auth-theme-switcher__choice'}
              onClick={() => setPreference(value)}>
              <Icon size={15}/><span>{label}</span>
            </button>
          ))}
        </div>
      </header>
      <div className="auth-stage">
        <SchoolIllustrationLeft />
        <section className="auth-panel" aria-label={title}>{children}</section>
        <SchoolIllustrationRight />
        <MobileIllustration />
      </div>
      <footer className="auth-footer"><span>© {new Date().getFullYear()} GrooveHub</span></footer>
    </main>
  );
};
