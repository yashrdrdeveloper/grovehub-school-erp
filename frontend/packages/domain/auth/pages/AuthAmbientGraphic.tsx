/**
 * Abstract academic pathways for the desktop sign-in panel.
 * Pure decorative SVG: no photographs, external assets, motion or duplicate RDR logo.
 */
export const AuthAmbientGraphic = () => (
  <div className="rdr-auth-art" aria-hidden="true">
    <svg className="rdr-auth-art__svg" viewBox="0 0 560 380" fill="none" xmlns="http://www.w3.org/2000/svg">
      <defs>
        <pattern id="rdr-auth-grid" width="28" height="28" patternUnits="userSpaceOnUse">
          <path d="M28 0H0V28" stroke="currentColor" strokeWidth=".6" opacity=".14"/>
        </pattern>
        <linearGradient id="rdr-auth-path" x1="44" y1="318" x2="496" y2="50" gradientUnits="userSpaceOnUse">
          <stop stopColor="#9CAFC2" stopOpacity=".13"/>
          <stop offset=".53" stopColor="#AF9469" stopOpacity=".75"/>
          <stop offset="1" stopColor="#9CAFC2" stopOpacity=".15"/>
        </linearGradient>
      </defs>
      <rect x="26" y="18" width="508" height="336" rx="25" fill="url(#rdr-auth-grid)"/>
      <rect x="50" y="43" width="460" height="284" rx="19" stroke="currentColor" strokeOpacity=".16" strokeWidth="1"/>
      <path d="M69 287C125 287 137 92 240 92S366 278 486 124" stroke="url(#rdr-auth-path)" strokeWidth="2.2"/>
      <path d="M69 227C151 227 174 314 280 314S392 99 489 99" stroke="currentColor" strokeOpacity=".21" strokeWidth="1.2" strokeDasharray="4 8"/>
      <path d="M74 164C183 164 207 68 307 68S423 243 487 243" stroke="currentColor" strokeOpacity=".16" strokeWidth="1.2"/>
      <path d="M160 92V52M408 279v38M279 160v-70M279 278v36" stroke="currentColor" strokeOpacity=".16" strokeDasharray="2 7"/>
      <rect x="190" y="128" width="180" height="130" rx="22" fill="currentColor" fillOpacity=".035" stroke="currentColor" strokeOpacity=".3"/>
      <rect x="206" y="143" width="148" height="100" rx="15" stroke="currentColor" strokeOpacity=".13"/>
      <path d="M240 175c16-5 27-3 40 5 13-8 24-10 40-5v37c-16-5-27-3-40 5-13-8-24-10-40-5v-37Z" stroke="currentColor" strokeWidth="2.6" strokeLinecap="round" strokeLinejoin="round"/>
      <path d="M280 180v37" stroke="currentColor" strokeWidth="2.4"/>
      <path d="M250 185c8-2 15-1 20 2M250 196c8-2 15-1 20 2M290 187c6-3 13-4 20-2M290 198c6-3 13-4 20-2" stroke="currentColor" strokeOpacity=".4" strokeWidth="1.5" strokeLinecap="round"/>
      <circle cx="160" cy="92" r="27" fill="var(--art-node-bg)" stroke="currentColor" strokeOpacity=".28"/>
      <path d="m149 95 11-6 11 6-11 6-11-6Zm0 0v11c7 5 15 5 22 0V95" stroke="currentColor" strokeWidth="2" strokeLinejoin="round"/>
      <circle cx="408" cy="279" r="27" fill="var(--art-node-bg)" stroke="currentColor" strokeOpacity=".28"/>
      <rect x="396" y="268" width="24" height="22" rx="3" stroke="currentColor" strokeWidth="2"/>
      <path d="M396 274h24M402 265v7M414 265v7M402 281h5M412 281h3" stroke="currentColor" strokeWidth="1.8" strokeLinecap="round"/>
      <circle cx="484" cy="123" r="5" fill="#B99B6B"/>
      <circle cx="75" cy="227" r="5" fill="#B99B6B"/>
      <path d="m478 56 2-8 2 8 8 2-8 2-2 8-2-8-8-2 8-2Z" fill="#B99B6B" fillOpacity=".72"/>
      <path d="m103 303 2-7 2 7 7 2-7 2-2 7-2-7-7-2 7-2Z" fill="currentColor" fillOpacity=".24"/>
    </svg>
    <div className="rdr-auth-art__caption">
      <span>LEARNING</span><span className="rdr-auth-art__sep"/> <span>PEOPLE</span><span className="rdr-auth-art__sep"/> <span>PROGRESS</span>
    </div>
  </div>
);
