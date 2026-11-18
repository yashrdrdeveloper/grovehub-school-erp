import type { ReactNode } from 'react';
import {
  ArrowRight, BookOpen, CalendarDays, GraduationCap, ListTodo, Monitor,
  Moon, School, Sparkles, Sun, UsersRound, WalletCards,
} from 'lucide-react';

import { SchoolMonogram } from '@Components';
import { useWorkspaceTheme, type ThemePreference } from '@Shell/theme';

const themeChoices: { value: ThemePreference; label: string; icon: typeof Sun }[] = [
  { value: 'system', label: 'System', icon: Monitor },
  { value: 'light', label: 'Light', icon: Sun },
  { value: 'dark', label: 'Dark', icon: Moon },
];

/** Decorative product journey; contains no personal data or live business metrics. */
const SchoolIllustration = () => (
  <div className="auth-art" aria-hidden="true">
    <div className="auth-art__quote">
      <span className="auth-art__quote-mark">“</span>
      <p>Everything your school needs, in one calm workspace.</p>
      <span className="auth-art__quote-underline" />
    </div>

    <div className="auth-art__feature auth-art__feature--crm">
      <div className="auth-art__feature-title">
        <UsersRound size={19} strokeWidth={1.6}/>
        <span>Lead CRM</span>
      </div>
      <div className="auth-art__divider"/>
      <div className="auth-art__journey">
        <span>Enquiry</span><ArrowRight size={12}/><span>Admission</span>
      </div>
      <div className="auth-art__contact"><i/><b/><span>✓</span></div>
      <div className="auth-art__contact"><i/><b/><span>✓</span></div>
    </div>

    <div className="auth-art__feature auth-art__feature--student">
      <div className="auth-art__feature-title">
        <GraduationCap size={19} strokeWidth={1.6}/>
        <span>Student Lifecycle 360°</span>
      </div>
      <small>Unified student view · UCV</small>
      <div className="auth-art__divider"/>
      <div className="auth-art__student">
        <span className="auth-art__student-avatar"/>
        <span className="auth-art__student-lines"><i/><i/></span>
        <span className="auth-art__student-ring"/>
      </div>
      <div className="auth-art__tags"><span>Academics</span><span>Growth</span></div>
    </div>

    <div className="auth-art__feature auth-art__feature--tasks">
      <ListTodo size={19} strokeWidth={1.6}/>
      <span><strong>Tasks</strong><small>Daily school work</small></span>
    </div>
    <div className="auth-art__feature auth-art__feature--finance">
      <WalletCards size={19} strokeWidth={1.6}/>
      <span><strong>Finance & Fees</strong><small>School finances</small></span>
    </div>

    <div className="auth-art__doodle auth-art__doodle--calendar"><CalendarDays size={25} strokeWidth={1.3}/></div>
    <div className="auth-art__doodle auth-art__doodle--book"><BookOpen size={27} strokeWidth={1.3}/></div>
    <div className="auth-art__doodle auth-art__doodle--school"><School size={23} strokeWidth={1.3}/></div>
    <svg className="auth-art__squiggle auth-art__squiggle--left" viewBox="0 0 170 60" fill="none">
      <path d="M3 39c22 33 33-29 50-18 18 13 0 38 23 28 22-10 22-47 44-30 17 13 19 24 47 13" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round"/>
    </svg>
    <svg className="auth-art__squiggle auth-art__squiggle--right" viewBox="0 0 170 60" fill="none">
      <path d="M3 38c22 30 38-31 55-19 18 13-3 39 22 26 21-11 29-43 47-25 15 15 24 20 40 13" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round"/>
    </svg>
    <div className="auth-art__dots auth-art__dots--left"/>
    <div className="auth-art__dots auth-art__dots--right"/>
    <div className="auth-art__spark auth-art__spark--left">✳</div>
    <div className="auth-art__spark auth-art__spark--right">✧</div>
    <div className="auth-art__floor"/>
    <div className="auth-art__motto">LEARN <span>·</span> ORGANIZE <span>·</span> GROW</div>
    <div className="auth-art__signature">RDR SCHOOL <span>/</span> GROOVEHUB</div>
    <div className="auth-art__mobile-story">
      <Sparkles size={14}/> From first enquiry to every milestone.
    </div>
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
        <SchoolIllustration/>
        <section className="auth-panel" aria-label={title}>
          {children}
        </section>
      </div>
      <footer className="auth-footer"><span>© {new Date().getFullYear()} GrooveHub</span></footer>
    </main>
  );
};
