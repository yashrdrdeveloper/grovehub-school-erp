import { Box, BOX_TAG } from '@Components/components/Box';
import type { ReactNode } from 'react';
import { renderAuthStatusBadge } from './renderer/authStatusArt.badge';
import type { AuthStatusKind } from './types/authStatusArt';
export type { AuthStatusKind } from './types/authStatusArt';
import { Check, Clock3, Link2Off, MailCheck, ShieldCheck, WifiOff } from '@Components/icons';

const statusIcons: Record<AuthStatusKind, ReactNode> = {
  email: <MailCheck size={52} strokeWidth={1.45} />,
  success: <Check size={50} strokeWidth={2.1} />,
  expired: <Link2Off size={46} strokeWidth={1.7} />,
  missing: <Link2Off size={46} strokeWidth={1.7} />,
  checking: <ShieldCheck size={48} strokeWidth={1.55} />,
  unavailable: <WifiOff size={46} strokeWidth={1.6} />,
  'attempts-exhausted': <Clock3 size={46} strokeWidth={1.6} />,
};

/** Reusable state-specific mobile illustration, independent of any auth API. */
export const AuthStatusArt = ({ kind }: { kind: AuthStatusKind }) => (
  <Box
    as={BOX_TAG.DIV}
    data-testid={`auth-status-${kind}`}
    className={`rdr-auth-mobile-status-art rdr-auth-mobile-status-art--${kind}`}
    aria-hidden="true"
  >
    <Box as={BOX_TAG.SPAN} className="rdr-auth-mobile-status-art__halo" />
    <Box as={BOX_TAG.SPAN} className="rdr-auth-mobile-status-art__symbol">
      {statusIcons[kind]}
    </Box>
    {renderAuthStatusBadge(kind)}
    <Box
      as={BOX_TAG.SPAN}
      className="rdr-auth-mobile-status-art__spark rdr-auth-mobile-status-art__spark--a"
    />
    <Box
      as={BOX_TAG.SPAN}
      className="rdr-auth-mobile-status-art__spark rdr-auth-mobile-status-art__spark--b"
    />
  </Box>
);
