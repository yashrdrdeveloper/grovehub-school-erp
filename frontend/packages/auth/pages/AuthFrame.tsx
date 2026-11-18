import type { ReactNode } from 'react';

// A single auth shell keeps branding, card dimensions and footer identical
// across login, password recovery and invitation setup.
export const AuthFrame = ({ children, title }: { children: ReactNode; title: string }) => (
  <main className="auth-page auth-page--consistent">
    <div className="auth-brand" aria-label="GrooveHub">
      <div className="brand-mark">G</div>
      <strong>groovehub<span className="brand-period">.</span></strong>
    </div>
    <section className="auth-panel" aria-label={title}>
      {children}
    </section>
    <footer className="auth-footer"><span>© {new Date().getFullYear()} GrooveHub</span></footer>
  </main>
);
