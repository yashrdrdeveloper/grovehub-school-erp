import type { ReactNode } from 'react';
import { ArrowRight, BookOpenCheck, GraduationCap, Megaphone, Monitor, Moon, School, Sparkles, Sun, WalletCards, Workflow, ListTodo, UsersRound } from 'lucide-react';

import { SchoolMonogram } from '@Components';
import { useWorkspaceTheme, type ThemePreference } from '@Shell/theme';

// Shared auth chrome: branding remains unchanged across all auth flows.
// The editorial illustration is decorative and desktop-only.
const themeChoices: { value: ThemePreference; label: string; icon: typeof Sun }[] = [
  { value: 'system', label: 'System', icon: Monitor },
  { value: 'light', label: 'Light', icon: Sun },
  { value: 'dark', label: 'Dark', icon: Moon },
];

export const AuthFrame = ({ children, title }: { children: ReactNode; title: string }) => {
  const { preference, setPreference } = useWorkspaceTheme();
  return (
  <main className="auth-page auth-page--consistent">
    <div className="auth-scene" aria-hidden="true">
      <div className="auth-scene__ground"/>
      <div className="auth-scene__prelude">
        <div className="auth-scene__prelude-illustration">
          <span className="auth-scene__prelude-trace"/>
          <span className="auth-scene__prelude-icon auth-scene__prelude-icon--one"><BookOpenCheck size={19} strokeWidth={1.6}/></span>
          <span className="auth-scene__prelude-icon auth-scene__prelude-icon--main"><School size={23} strokeWidth={1.55}/></span>
          <span className="auth-scene__prelude-icon auth-scene__prelude-icon--three"><GraduationCap size={20} strokeWidth={1.6}/></span>
        </div>
        <p>One connected workspace for every school day.</p>
      </div>
      <div className="auth-scene__sheet auth-scene__sheet--left">
        <div className="auth-scene__sheet-head">
          <UsersRound size={23} strokeWidth={1.6}/>
          <span>Lead CRM</span>
        </div>
        <div className="auth-scene__sheet-rule"/>
        <div className="auth-scene__journey">
          <span>Enquiry</span><ArrowRight size={13}/><span>Admission</span>
        </div>
        <div className="auth-scene__sheet-entry"><span className="auth-scene__sheet-avatar"/><span className="auth-scene__sheet-line"/><span className="auth-scene__sheet-check">✓</span></div>
        <div className="auth-scene__sheet-entry"><span className="auth-scene__sheet-avatar"/><span className="auth-scene__sheet-line auth-scene__sheet-line--short"/><span className="auth-scene__sheet-check">✓</span></div>
      </div>
      <div className="auth-scene__sheet auth-scene__sheet--right">
        <div className="auth-scene__sheet-head">
          <UsersRound size={23} strokeWidth={1.6}/>
          <span>Student Lifecycle 360°</span>
        </div>
        <div className="auth-scene__sheet-subtitle">Unified student view · UCV</div>
        <div className="auth-scene__sheet-rule"/>
        <div className="auth-scene__student-profile">
          <span className="auth-scene__student-avatar"/>
          <span className="auth-scene__student-lines"><i/><i/></span>
          <span className="auth-scene__student-ring"/>
        </div>
        <div className="auth-scene__student-tags"><span>Academics</span><span>Growth</span></div>
      </div>
      <div className="auth-scene__dots auth-scene__dots--left"/>
      <div className="auth-scene__dots auth-scene__dots--right"/>
      <svg className="auth-scene__squiggle auth-scene__squiggle--left" viewBox="0 0 190 74" fill="none"><path d="M3 49c20 36 39-36 50-22 15 18-16 40 6 36 22-5 22-51 43-42 19 9 17 47 42 29 19-13 26-16 43-12" stroke="currentColor" strokeWidth="1.3" strokeLinecap="round"/></svg>
      <svg className="auth-scene__squiggle auth-scene__squiggle--right" viewBox="0 0 170 64" fill="none"><path d="M3 40c19 27 32-24 49-21 15 3-1 31 15 31 20 0 25-46 47-38 16 6 9 41 29 32 8-4 16-9 24-6" stroke="currentColor" strokeWidth="1.3" strokeLinecap="round"/></svg>
      <div className="auth-scene__orbit auth-scene__orbit--one"/>
      <div className="auth-scene__orbit auth-scene__orbit--two"/>
      <div className="auth-scene__badge auth-scene__badge--calendar">
        <svg viewBox="0 0 50 50" fill="none" aria-hidden="true"><rect x="8" y="11" width="34" height="31" rx="4" stroke="currentColor" strokeWidth="1.6"/><path d="M8 20h34M17 7v9M33 7v9M17 28h5M28 28h5M17 35h5" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round"/></svg>
      </div>
      <div className="auth-scene__badge auth-scene__badge--book">
        <svg viewBox="0 0 58 48" fill="none" aria-hidden="true"><path d="M29 40c-8-6-17-7-25-5V8c9-2 18 0 25 6 8-6 17-8 25-6v27c-8-2-17-1-25 5Z" stroke="currentColor" strokeWidth="1.6" strokeLinejoin="round"/><path d="M29 14v26M11 16c5-1 9 0 13 2M35 18c5-2 9-3 13-2M11 23c5-1 9 0 13 2M35 25c5-2 9-3 13-2" stroke="currentColor" strokeWidth="1.3" strokeLinecap="round"/></svg>
      </div>
      <div className="auth-scene__badge auth-scene__badge--cap">
        <svg viewBox="0 0 58 50" fill="none" aria-hidden="true"><path d="M5 21 29 10l24 11-24 11L5 21Z" stroke="currentColor" strokeWidth="1.7" strokeLinejoin="round"/><path d="M15 26v11c10 8 20 8 29 0V26M53 21v19" stroke="currentColor" strokeWidth="1.7" strokeLinecap="round"/><circle cx="53" cy="43" r="3" fill="currentColor"/></svg>
      </div>
      <div className="auth-scene__badge auth-scene__badge--pencil">
        <svg viewBox="0 0 60 60" fill="none" aria-hidden="true"><path d="m12 44 4-14L42 4l14 14-26 26-14 4-4-4Z" stroke="currentColor" strokeWidth="1.6" strokeLinejoin="round"/><path d="m35 11 14 14M16 30l14 14M12 44l4 4" stroke="currentColor" strokeWidth="1.5"/></svg>
      </div>
      <svg className="auth-scene__paperplane" viewBox="0 0 120 75" fill="none" aria-hidden="true">
        <path d="m6 31 101-27-36 61-17-23L6 31Zm48 11 53-38-34 43" stroke="currentColor" strokeWidth="1.5" strokeLinejoin="round"/>
        <path d="M5 62c15-11 24-6 35-11" stroke="currentColor" strokeWidth="1.3" strokeDasharray="3 6" strokeLinecap="round"/>
      </svg>
      <div className="auth-scene__spark auth-scene__spark--left">✳</div>
      <div className="auth-scene__spark auth-scene__spark--right">✧</div>
      <div className="auth-scene__near auth-scene__near--left">
        <div className="auth-scene__near-icon"><ListTodo size={17}/></div>
        <div><strong>Tasks</strong><span>Stay on top of daily work.</span></div>
      </div>
      <div className="auth-scene__near auth-scene__near--right">
        <div className="auth-scene__near-icon"><WalletCards size={17}/></div>
        <div><strong>Finance & fees</strong><span>Keep finances organized.</span></div>
      </div>
      <div className="auth-scene__quote">
        <span className="auth-scene__quote-mark">“</span>
        <p>Everything your school needs, in one calm workspace.</p>
        <span className="auth-scene__quote-rule"/>
      </div>
      <div className="auth-scene__features">
        <span><Megaphone size={14}/> Campaigns</span>
        <span><Workflow size={14}/> Automation</span>
        <span><Sparkles size={14}/> AI assistance</span>
      </div>
      <div className="auth-scene__micro auth-scene__micro--left">LEARN <span>·</span> ORGANIZE <span>·</span> GROW</div>
      <div className="auth-scene__micro auth-scene__micro--right">RDR SCHOOL <span> / </span> GROOVEHUB</div>
    </div>
    <header className="auth-topbar">
      <div className="auth-brand auth-brand--school" aria-label="RDR School">
        <SchoolMonogram />
        <strong>RDR School</strong>
      </div>
      <div className="auth-theme-switcher" role="group" aria-label="Appearance">
        {themeChoices.map(({ value, label, icon: Icon }) => (
          <button key={value} type="button" title={label} aria-label={`${label} theme`} aria-pressed={preference === value}
            className={preference === value ? 'auth-theme-switcher__choice is-active' : 'auth-theme-switcher__choice'}
            onClick={() => setPreference(value)}>
            <Icon size={15}/><span>{label}</span>
          </button>
        ))}
      </div>
    </header>
    <div className="auth-compact-art" aria-hidden="true">
      <svg className="auth-compact-art__sketch" viewBox="0 0 480 250" preserveAspectRatio="none" fill="none">
        <path d="M16 72c19 31 35-29 54-10 15 15 12 31 27 17M378 57c24-25 27 24 48 8 13-11 22-17 38-9"
          stroke="currentColor" strokeWidth="1.5" strokeLinecap="round"/>
        <path d="M65 194c42-23 82-17 129-6 52 12 73 13 122-2 38-11 69-7 100 7"
          stroke="currentColor" strokeWidth="1.2" strokeDasharray="4 7" strokeLinecap="round" opacity=".65"/>
        <circle cx="29" cy="148" r="4" stroke="currentColor"/>
        <circle cx="443" cy="129" r="4" stroke="currentColor"/>
      </svg>
      <div className="auth-compact-art__spark auth-compact-art__spark--left">✳</div>
      <div className="auth-compact-art__spark auth-compact-art__spark--right">✦</div>
      <div className="auth-compact-art__item auth-compact-art__item--crm">
        <div className="auth-compact-art__card-top">
          <span className="auth-compact-art__icon"><UsersRound size={17} strokeWidth={1.7}/></span>
          <span className="auth-compact-art__copy"><strong>Lead CRM</strong><small>Admissions journey</small></span>
        </div>
        <div className="auth-compact-art__journey"><span>Enquiry</span><ArrowRight size={12}/><span>Admission</span></div>
        <div className="auth-compact-art__lines"><i/><i/></div>
      </div>
      <div className="auth-compact-art__item auth-compact-art__item--student">
        <div className="auth-compact-art__card-top">
          <span className="auth-compact-art__icon"><GraduationCap size={18} strokeWidth={1.7}/></span>
          <span className="auth-compact-art__copy"><strong>Student 360°</strong><small>Lifecycle · UCV</small></span>
        </div>
        <div className="auth-compact-art__student-preview">
          <span className="auth-compact-art__avatar"/>
          <span className="auth-compact-art__lines"><i/><i/></span>
          <span className="auth-compact-art__progress"/>
        </div>
      </div>
      <div className="auth-compact-art__story">
        <span>“</span>
        <p>From first enquiry to every milestone.</p>
        <small>LEARN <b>·</b> ORGANIZE <b>·</b> GROW</small>
      </div>
      <div className="auth-compact-art__baseline"/>
    </div>
    <section className="auth-panel" aria-label={title}>
      {children}
    </section>
    <footer className="auth-footer"><span>© {new Date().getFullYear()} GrooveHub</span></footer>
  </main>
  );
};
