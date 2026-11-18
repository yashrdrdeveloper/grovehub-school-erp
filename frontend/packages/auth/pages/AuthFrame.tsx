import type { ReactNode } from 'react';

// One shared auth frame for login, recovery and invitation setup.
// Illustrations are decorative, hidden from assistive technology and mobile screens.
export const AuthFrame = ({ children, title }: { children: ReactNode; title: string }) => (
  <main className="auth-page auth-page--consistent">
    <div className="auth-school-art" aria-hidden="true">
      <div className="auth-school-art__halo auth-school-art__halo--left"/>
      <div className="auth-school-art__halo auth-school-art__halo--right"/>
      <svg className="auth-school-art__campus" viewBox="0 0 360 440" fill="none" xmlns="http://www.w3.org/2000/svg">
        <circle cx="172" cy="222" r="141" fill="#F9DFB8" opacity=".55"/>
        <circle cx="62" cy="83" r="6" fill="#EBAC75"/>
        <path d="M38 112h30M53 97v30" stroke="#A98673" strokeWidth="2" strokeLinecap="round"/>
        <path d="M28 341h302" stroke="#32333B" strokeWidth="2" strokeLinecap="round"/>
        <rect x="70" y="189" width="207" height="150" rx="3" fill="#FFFDF9" stroke="#2D3039" strokeWidth="2.5"/>
        <rect x="93" y="150" width="160" height="42" fill="#F8F0E7" stroke="#2D3039" strokeWidth="2.5"/>
        <path d="M81 151l91-57 94 57H81Z" fill="#F1AA91" stroke="#2D3039" strokeWidth="2.5" strokeLinejoin="round"/>
        <path d="M119 150v-20m106 20v-20" stroke="#2D3039" strokeWidth="2"/>
        <circle cx="173" cy="136" r="14" fill="#FFF9EE" stroke="#2D3039" strokeWidth="2"/>
        <path d="M173 128v9l6 4" stroke="#2D3039" strokeWidth="2" strokeLinecap="round"/>
        <rect x="143" y="263" width="59" height="76" rx="26" fill="#CFDCE3" stroke="#2D3039" strokeWidth="2.5"/>
        <path d="M172 264v75" stroke="#2D3039" strokeWidth="2"/>
        <rect x="93" y="215" width="31" height="35" rx="2" fill="#F5CE81" stroke="#2D3039" strokeWidth="2"/>
        <path d="M108 215v35m-15-17h31" stroke="#2D3039" strokeWidth="1.5"/>
        <rect x="219" y="215" width="32" height="35" rx="2" fill="#F5CE81" stroke="#2D3039" strokeWidth="2"/>
        <path d="M235 215v35m-16-17h32" stroke="#2D3039" strokeWidth="1.5"/>
        <rect x="93" y="269" width="31" height="35" rx="2" fill="#F5CE81" stroke="#2D3039" strokeWidth="2"/>
        <path d="M108 269v35m-15-17h31" stroke="#2D3039" strokeWidth="1.5"/>
        <rect x="219" y="269" width="32" height="35" rx="2" fill="#F5CE81" stroke="#2D3039" strokeWidth="2"/>
        <path d="M235 269v35m-16-17h32" stroke="#2D3039" strokeWidth="1.5"/>
        <path d="M172 94V66m0 0 29 11-29 11" fill="#EEAA8A" stroke="#2D3039" strokeWidth="2" strokeLinejoin="round"/>
        <path d="M20 342c15-25 17-53 12-76m-12 76c-15-26-12-43-5-54" stroke="#3F776D" strokeWidth="4" strokeLinecap="round"/>
        <circle cx="29" cy="260" r="26" fill="#A3C9B2" stroke="#2D3039" strokeWidth="2"/>
        <circle cx="14" cy="285" r="20" fill="#A3C9B2" stroke="#2D3039" strokeWidth="2"/>
        <circle cx="316" cy="282" r="30" fill="#A3C9B2" stroke="#2D3039" strokeWidth="2"/>
        <circle cx="331" cy="308" r="18" fill="#A3C9B2" stroke="#2D3039" strokeWidth="2"/>
        <path d="M319 309v33" stroke="#32685F" strokeWidth="5"/>
        <rect x="107" y="366" width="132" height="35" rx="17.5" fill="#FFFDF9" stroke="#2D3039" strokeWidth="1.5"/>
        <path d="M123 383h13m-7-6v12" stroke="#DE8B79" strokeWidth="2" strokeLinecap="round"/>
        <text x="146" y="388" fill="#33313B" fontSize="12" fontFamily="sans-serif" fontWeight="600">RDR SCHOOL</text>
      </svg>
      <svg className="auth-school-art__learning" viewBox="0 0 370 440" fill="none" xmlns="http://www.w3.org/2000/svg">
        <circle cx="176" cy="227" r="143" fill="#E4E9E7" opacity=".9"/>
        <circle cx="272" cy="108" r="32" fill="#F4D79F"/>
        <path d="M273 90v36m-18-18h36" stroke="#5B5C68" strokeWidth="1.6" strokeLinecap="round"/>
        <rect x="27" y="102" width="118" height="90" rx="9" fill="#FFFDF9" stroke="#30343B" strokeWidth="2" transform="rotate(-7 27 102)"/>
        <path d="M47 125h51m-51 15h69m-69 15h44" stroke="#8E9E9D" strokeWidth="4" strokeLinecap="round"/>
        <rect x="96" y="156" width="215" height="178" rx="10" fill="#F4C9B4" stroke="#30343B" strokeWidth="2.5"/>
        <path d="M119 193c20-11 47-12 69-1v116c-22-11-49-10-69 1V193Zm69-1c22-11 52-11 79 1v116c-27-12-57-11-79-1V192Z" fill="#FFFDF9" stroke="#30343B" strokeWidth="2.5" strokeLinejoin="round"/>
        <path d="M135 220h36m-36 16h36m-36 16h29m42-32h43m-43 16h43m-43 16h32" stroke="#B3B2AF" strokeWidth="2" strokeLinecap="round"/>
        <path d="M188 191v117" stroke="#30343B" strokeWidth="2.5"/>
        <path d="M90 335h232" stroke="#30343B" strokeWidth="2.5" strokeLinecap="round"/>
        <path d="M168 102l42-20 43 20-43 18-42-18Z" fill="#34343E" stroke="#34343E" strokeWidth="2" strokeLinejoin="round"/>
        <path d="M183 110v19c17 16 39 16 55 0v-19" fill="#6A687A" stroke="#34343E" strokeWidth="2"/>
        <path d="M253 105v31" stroke="#34343E" strokeWidth="2"/>
        <circle cx="253" cy="137" r="4" fill="#F2B76E"/>
        <rect x="258" y="352" width="66" height="24" rx="12" fill="#FFFDF9" stroke="#30343B" strokeWidth="1.5"/>
        <circle cx="276" cy="364" r="5" fill="#71B7A1"/>
        <path d="M286 364h24" stroke="#A4A6A7" strokeWidth="2" strokeLinecap="round"/>
        <path d="M36 333c11-10 24-10 35 0 11-10 24-10 35 0" stroke="#3F776D" strokeWidth="3" strokeLinecap="round"/>
        <path d="M22 70h17m-8-9v18" stroke="#D59A6D" strokeWidth="2" strokeLinecap="round"/>
      </svg>
      <div className="auth-school-art__caption auth-school-art__caption--left">A place to learn.</div>
      <div className="auth-school-art__caption auth-school-art__caption--right">A smarter way to manage.</div>
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
