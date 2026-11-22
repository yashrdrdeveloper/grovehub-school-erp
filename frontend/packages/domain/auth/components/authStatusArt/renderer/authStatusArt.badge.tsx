import { Box, BOX_TAG } from '@Components/components/Box';
import type { ReactNode } from 'react';
import { CheckCircle2 } from '@Components/icons';
import type { AuthStatusKind } from '../types/authStatusArt';

/** JSX is selected outside the parent presentation, with semantic theme tokens. */
export const renderAuthStatusBadge = (kind: AuthStatusKind): ReactNode => {
  if (kind === 'email' || kind === 'success') {
    return (
      <CheckCircle2
        className="rdr-auth-mobile-status-art__badge"
        size={32}
        fill="var(--gh-color-success)"
        color="var(--gh-color-primary-foreground)"
        strokeWidth={2.5}
      />
    );
  }
  if (kind === 'expired' || kind === 'missing') {
    return (
      <Box as={BOX_TAG.SPAN} className="rdr-auth-mobile-status-art__x">
        ×
      </Box>
    );
  }
  return null;
};
