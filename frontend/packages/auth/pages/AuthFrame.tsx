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
      <div className="auth-scene__micro auth-scene__micro--left">LEARN <span>·</span> ORGANIZE <span>·</span> GROW</div>
      <div className="auth-scene__micro auth-scene__micro--right">RDR SCHOOL <span> / </span> GROOVEHUB</div>
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
