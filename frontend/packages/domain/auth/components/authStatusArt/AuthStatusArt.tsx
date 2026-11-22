import type { ReactNode } from 'react';
import { Check, CheckCircle2, Clock3, Link2Off, MailCheck, ShieldCheck, WifiOff } from '@Components/icons';

export type AuthStatusKind = 'email' | 'success' | 'expired' | 'missing' | 'checking' | 'unavailable' | 'attempts-exhausted';

const statusIcons: Record<AuthStatusKind, ReactNode> = {
  email: <MailCheck size={52} strokeWidth={1.45}/>,
  success: <Check size={50} strokeWidth={2.1}/>,
  expired: <Link2Off size={46} strokeWidth={1.7}/>,
  missing: <Link2Off size={46} strokeWidth={1.7}/>,
  checking: <ShieldCheck size={48} strokeWidth={1.55}/>,
  unavailable: <WifiOff size={46} strokeWidth={1.6}/>,
  'attempts-exhausted': <Clock3 size={46} strokeWidth={1.6}/>,
};

/** Reusable state-specific mobile illustration, independent of any auth API. */
export const AuthStatusArt = ({ kind }: { kind: AuthStatusKind }) => (
  <div data-testid={`auth-status-${kind}`} className={`rdr-auth-mobile-status-art rdr-auth-mobile-status-art--${kind}`} aria-hidden="true">
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

