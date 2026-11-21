import type { ReactNode } from 'react';
import { Check, CheckCircle2, Clock3, Link2Off, MailCheck, ShieldCheck, WifiOff } from 'lucide-react';

type StatusKind = 'email' | 'success' | 'expired' | 'missing' | 'checking' | 'unavailable' | 'attempts-exhausted';

/**
 * The mobile auth artwork is an illustration, not a rasterized screenshot.
 * It stays crisp at every viewport width and is self-contained/offline.
 */
export const AuthMobileSchoolScene = () => (
  <svg className="rdr-auth-mobile-scene" viewBox="0 0 393 250" fill="none" xmlns="http://www.w3.org/2000/svg"
    preserveAspectRatio="xMidYMax slice" focusable="false" aria-hidden="true">
    <defs>
      <linearGradient id="auth-scene-school" x1="291" y1="105" x2="379" y2="199" gradientUnits="userSpaceOnUse">
        <stop stopColor="#D2F5FF"/><stop offset=".52" stopColor="#5DBBDD"/><stop offset="1" stopColor="#08749D"/>
      </linearGradient>
      <linearGradient id="auth-scene-lawn" x1="120" y1="215" x2="351" y2="234" gradientUnits="userSpaceOnUse">
        <stop stopColor="#076B92"/><stop offset=".58" stopColor="#0796A6"/><stop offset="1" stopColor="#2DD3C4"/>
      </linearGradient>
      <linearGradient id="auth-scene-bank" x1="209" y1="208" x2="377" y2="178" gradientUnits="userSpaceOnUse">
        <stop stopColor="#0084A5"/><stop offset="1" stopColor="#3CE1CB"/>
      </linearGradient>
    </defs>
    <path d="M0 192C57 175 92 191 136 178S228 174 274 187 347 179 393 161V250H0Z" fill="#063768" opacity=".72"/>
    <path d="M0 217C65 189 98 205 146 190S230 191 281 197 352 177 393 174V250H0Z" fill="#085987" opacity=".78"/>
    <g className="rdr-auth-mobile-scene__school" transform="translate(-150 -50) scale(1.35)">
      <ellipse cx="343" cy="200" rx="65" ry="8" fill="#013D63" opacity=".3"/>
      <path d="M339 120L390 130V184L339 187Z" fill="#126E9A" stroke="#77CDE4" strokeWidth="1.5"/>
      <path d="M339 120L353 105 393 117V132Z" fill="#0A5986" stroke="#91E3F0" strokeWidth="1.5"/>
      <path d="M353 139L376 137V166L353 170Z" fill="#0A4771"/>
      <path d="M356 143L372 141M356 151L372 149M356 159L372 157" stroke="#96E7F1" strokeWidth="2"/>
      <path d="M280 139L310 127V184H280Z" fill="#227FA8" stroke="#72CDE4" strokeWidth="1.5"/>
      <path d="M276 137L300 123 316 125 292 142Z" fill="#0D608E" stroke="#A8EDFB" strokeWidth="2"/>
      <path d="M284 147H303V175H284Z" fill="#115580"/>
      <path d="M287 151H300M287 159H300M287 167H300" stroke="#8DE0F0" strokeWidth="2.3"/>
      <path d="M304 124L343 88 380 121 376 131 306 132Z" fill="#74C9E1" stroke="#E1FBFF" strokeWidth="2.6"/>
      <path d="M312 124L343 96 372 124Z" fill="#0B729D" stroke="#C5F7FF" strokeWidth="2"/>
      <path d="M316 133V187H370V133Z" fill="url(#auth-scene-school)" stroke="#8EDFF0" strokeWidth="1.5"/>
      <path d="M314 132H373M311 187H375M309 192H378" stroke="#D9F9FF" strokeWidth="3"/>
      <path d="M324 137V182M336 137V182M349 137V182M362 137V182" stroke="#ECFFFF" strokeWidth="5.4"/>
      <path d="M316 181H370" stroke="#74CCDE" strokeWidth="3"/>
      <path d="M340 187V158C340 149 348 149 348 158V187Z" fill="#05466E" stroke="#B8EFF8" strokeWidth="1.6"/>
      <circle cx="343" cy="115" r="6" fill="#E8FEFF" stroke="#72C3E1" strokeWidth="1.5"/>
      <path d="M343 111V115L346 117" stroke="#287499" strokeWidth="1.5" strokeLinecap="round"/>
      <path d="M344 88V64M344 65L363 70 344 76Z" stroke="#E1FAFF" strokeWidth="2.2" strokeLinejoin="round"/>
      <path d="M376 132L393 128" stroke="#E0FAFF" strokeWidth="2"/>
    </g>
    <g className="rdr-auth-mobile-scene__trees">
      <circle cx="279" cy="184" r="14" fill="#07688D"/><circle cx="268" cy="194" r="13" fill="#067DA0"/>
      <circle cx="291" cy="195" r="15" fill="#0D85A1"/>
      <circle cx="379" cy="181" r="17" fill="#0A9EAA"/>
      <circle cx="393" cy="188" r="16" fill="#087C9D"/>
    </g>
    <path d="M0 226C83 197 114 235 173 205S270 199 310 204 357 190 393 182V250H0Z" fill="url(#auth-scene-lawn)"/>
    <path d="M0 242C70 212 134 242 200 220S322 208 393 196V250H0Z" fill="url(#auth-scene-bank)" opacity=".88"/>
    <path d="M0 252C108 218 162 254 255 227S350 214 393 206V250H0Z" fill="#1AB6B6" opacity=".82"/>
  </svg>
);

/** A deliberately simpler, sharper logo for mobile, matching the approved RDR letters. */
export const AuthMobileRdrMark = () => (
  <span className="rdr-auth-mobile-mark" aria-hidden="true">
    <svg viewBox="0 0 72 72" xmlns="http://www.w3.org/2000/svg" focusable="false">
      <rect width="72" height="72" rx="17" fill="#FFFFFF"/>
      <text x="36" y="44" textAnchor="middle" fill="#122B55"
        fontFamily="Arial, Helvetica, sans-serif" fontWeight="900" fontSize="28" letterSpacing="-2.9">RDR</text>
      <path d="M18 49C29 45 44 45 55 48" fill="none" stroke="#04B9B4" strokeWidth="2.6" strokeLinecap="round"/>
    </svg>
  </span>
);

const statusIcons: Record<StatusKind, ReactNode> = {
  email: <MailCheck size={52} strokeWidth={1.45}/>,
  success: <Check size={50} strokeWidth={2.1}/>,
  expired: <Link2Off size={46} strokeWidth={1.7}/>,
  missing: <Link2Off size={46} strokeWidth={1.7}/>,
  checking: <ShieldCheck size={48} strokeWidth={1.55}/>,
  unavailable: <WifiOff size={46} strokeWidth={1.6}/>,
  'attempts-exhausted': <Clock3 size={46} strokeWidth={1.6}/>,
};

/** Reusable state-specific mobile illustration, independent of any auth API. */
export const AuthMobileStatusArt = ({ kind }: { kind: StatusKind }) => (
  <div className={`rdr-auth-mobile-status-art rdr-auth-mobile-status-art--${kind}`} aria-hidden="true">
    <span className="rdr-auth-mobile-status-art__halo"/>
    <span className="rdr-auth-mobile-status-art__symbol">{statusIcons[kind]}</span>
    {kind === 'email' || kind === 'success'
      ? <CheckCircle2 className="rdr-auth-mobile-status-art__badge" size={32} fill="#039C90" color="#fff" strokeWidth={2.5}/>
      : kind === 'expired' || kind === 'missing'
        ? <span className="rdr-auth-mobile-status-art__x">×</span>
        : null}
    <span className="rdr-auth-mobile-status-art__spark rdr-auth-mobile-status-art__spark--a"/>
    <span className="rdr-auth-mobile-status-art__spark rdr-auth-mobile-status-art__spark--b"/>
  </div>
);
