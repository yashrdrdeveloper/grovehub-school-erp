import { Box, BOX_TAG } from '@Components/components/Box';
import styles from './AuthStatusArt.module.css';
import type { ReactNode } from 'react';
import { AUTH_STATUS_KIND } from './constants/authStatusArt.kind';
import { renderAuthStatusBadge } from './renderer/authStatusArt.badge';
import type { AuthStatusKind } from './types/authStatusArt';
export type { AuthStatusKind } from './types/authStatusArt';
import { Check, Clock3, Link2Off, MailCheck, ShieldCheck, WifiOff } from '@Components/icons';

const statusIcons: Record<AuthStatusKind, ReactNode> = {
  [AUTH_STATUS_KIND.EMAIL]: <MailCheck size={52} strokeWidth={1.45} />,
  [AUTH_STATUS_KIND.SUCCESS]: <Check size={50} strokeWidth={2.1} />,
  [AUTH_STATUS_KIND.EXPIRED]: <Link2Off size={46} strokeWidth={1.7} />,
  [AUTH_STATUS_KIND.MISSING]: <Link2Off size={46} strokeWidth={1.7} />,
  [AUTH_STATUS_KIND.CHECKING]: <ShieldCheck size={48} strokeWidth={1.55} />,
  [AUTH_STATUS_KIND.UNAVAILABLE]: <WifiOff size={46} strokeWidth={1.6} />,
  [AUTH_STATUS_KIND.ATTEMPTS_EXHAUSTED]: <Clock3 size={46} strokeWidth={1.6} />,
};

/** Reusable state-specific mobile illustration, independent of any auth API. */
export const AuthStatusArt = ({ kind }: { kind: AuthStatusKind }) => (
  <Box
    as={BOX_TAG.DIV}
    data-testid={`auth-status-${kind}`}
    className={styles['rdr-auth-mobile-status-art']}
    aria-hidden="true"
  >
    <Box as={BOX_TAG.SPAN} className={styles['rdr-auth-mobile-status-art__halo']} />
    <Box as={BOX_TAG.SPAN} className={styles['rdr-auth-mobile-status-art__symbol']}>
      {statusIcons[kind]}
    </Box>
    {renderAuthStatusBadge(kind)}
    <Box
      as={BOX_TAG.SPAN}
      className={`${styles['rdr-auth-mobile-status-art__spark']} ${styles['rdr-auth-mobile-status-art__spark--a']}`}
    />
    <Box
      as={BOX_TAG.SPAN}
      className={`${styles['rdr-auth-mobile-status-art__spark']} ${styles['rdr-auth-mobile-status-art__spark--b']}`}
    />
  </Box>
);
