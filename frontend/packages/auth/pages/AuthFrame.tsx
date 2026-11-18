import type { ReactNode } from 'react';
import { ArrowRight, Monitor, Moon, Sun } from 'lucide-react';

import { SchoolMonogram } from '@Components';
import { useWorkspaceTheme, type ThemePreference } from '@Shell/theme';

const themeChoices: { value: ThemePreference; label: string; icon: typeof Sun }[] = [
  { value: 'system', label: 'System', icon: Monitor },
  { value: 'light', label: 'Light', icon: Sun },
  { value: 'dark', label: 'Dark', icon: Moon },
];

/** The illustrations are decorative product concepts, not live school data. */
const SchoolProductScene = () => (
  <div className="auth-product-scene" aria-hidden="true">
    <div className="auth-product-scene__headline">
      <span className="auth-product-scene__eyebrow">ONE CONNECTED SCHOOL</span>
      <p>From first enquiry to every milestone.</p>
    </div>

    <div className="auth-product-card auth-product-card--crm">
      <div className="auth-product-card__header">
        <span className="auth-product-card__symbol">◎</span>
        <div><strong>Lead CRM</strong><small>Admissions journey</small></div>
      </div>
      <div className="auth-product-card__flow">
        <span>Enquiry</span><ArrowRight size={13}/><span>Admission</span>
      </div>
      <div className="auth-product-card__line auth-product-card__line--long"/>
      <div className="auth-product-card__line auth-product-card__line--short"/>
    </div>

    <div className="auth-product-card auth-product-card--student">
      <div className="auth-product-card__header">
        <span className="auth-product-card__symbol">◉</span>
        <div><strong>Student Lifecycle 360°</strong><small>Unified student view · UCV</small></div>
      </div>
      <div className="auth-product-card__profile">
        <div className="auth-product-card__avatar"/>
        <div className="auth-product-card__profile-lines"><i/><i/></div>
        <span className="auth-product-card__profile-ring"/>
      </div>
      <div className="auth-product-card__stages"><span>Academics</span><span>Wellbeing</span><span>Growth</span></div>
    </div>
  </div>
);

export const AuthFrame = ({ children, title }: { children: ReactNode; title: string }) => {
  const { preference, setPreference } = useWorkspaceTheme();

  return (
    <main className="auth-page auth-page--consistent">
      <SchoolProductScene />
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
      <section className="auth-panel" aria-label={title}>
        {children}
      </section>
      <footer className="auth-footer"><span>© {new Date().getFullYear()} GrooveHub</span></footer>
    </main>
  );
};
