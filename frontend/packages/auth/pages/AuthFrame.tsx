import type { ReactNode } from 'react';

// Shared auth chrome: branding remains unchanged across all auth flows.
// The editorial illustration is decorative and desktop-only.
export const AuthFrame = ({ children, title }: { children: ReactNode; title: string }) => (
  <main className="auth-page auth-page--consistent">
    <div className="auth-scene" aria-hidden="true">
      <div className="auth-scene__ground"/>
      <div className="auth-scene__note auth-scene__note--left"><i/><i/><i/></div>
      <div className="auth-scene__note auth-scene__note--right"><i/><i/></div>
      <div className="auth-scene__dots auth-scene__dots--left"/>
      <div className="auth-scene__dots auth-scene__dots--right"/>
      <svg className="auth-scene__squiggle auth-scene__squiggle--left" viewBox="0 0 190 74" fill="none"><path d="M3 49c20 36 39-36 50-22 15 18-16 40 6 36 22-5 22-51 43-42 19 9 17 47 42 29 19-13 26-16 43-12" stroke="currentColor" strokeWidth="1.3" strokeLinecap="round"/></svg>
      <svg className="auth-scene__squiggle auth-scene__squiggle--right" viewBox="0 0 170 64" fill="none"><path d="M3 40c19 27 32-24 49-21 15 3-1 31 15 31 20 0 25-46 47-38 16 6 9 41 29 32 8-4 16-9 24-6" stroke="currentColor" strokeWidth="1.3" strokeLinecap="round"/></svg>
      <div className="auth-scene__orbit auth-scene__orbit--one"/>
      <div className="auth-scene__orbit auth-scene__orbit--two"/>
    </div>
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
