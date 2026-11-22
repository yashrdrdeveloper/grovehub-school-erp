import { Button } from '@Components/components/Button';
import { Box } from '@Components/components/Box';
import { ArrowLeft } from '@Components/icons';
import { Link } from 'react-router-dom';
import { AUTH_ROUTE } from '@Business/auth/routes';

import { RecoveryArtwork } from '../recoveryArtwork';
import type { RecoveryStateComponent } from './types/recoveryState';

/** Non-form recovery states, retry and accessible next actions. */
export const RecoveryState: RecoveryStateComponent = ({ view, purpose, invitation, slow, onRetry }) => {
  let checkingMessage = 'Verifying your reset link securely…';
  if (invitation) checkingMessage = 'Verifying your invitation securely…';
  if (slow) checkingMessage = 'Still checking. This may take a moment…';
  const showResetLink = !invitation && (view === 'missing' || view === 'expired');
  const showBack = view === 'missing' || view === 'expired' || view === 'attempts-exhausted';

  return (
    <Box as="div" className="rdr-recovery__state">
      <RecoveryArtwork view={view} purpose={purpose}/>
      {view === 'checking' && <Box as="p" className="rdr-recovery__slow" role="status">{checkingMessage}</Box>}
      {view === 'success' && <Link className="rdr-login-v2__submit rdr-recovery__state-button" to={AUTH_ROUTE.LOGIN}>Go to login</Link>}
      {view === 'unavailable' && <Button unstyled className="rdr-login-v2__submit rdr-recovery__state-button"
        type="button" onClick={onRetry}>Try checking again</Button>}
      {showResetLink && <Link className="rdr-login-v2__submit rdr-recovery__state-button"
        to={AUTH_ROUTE.FORGOT_PASSWORD}>Request new link</Link>}
      {showBack && <Link className="rdr-recovery__back" to={AUTH_ROUTE.LOGIN}>
        <ArrowLeft size={16}/> Return to login
      </Link>}
    </Box>
  );
};
