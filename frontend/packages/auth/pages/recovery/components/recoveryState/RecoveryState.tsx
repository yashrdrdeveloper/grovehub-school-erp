import { ArrowLeft } from 'lucide-react';
import { Link } from 'react-router-dom';

import { RecoveryArtwork } from '../recoveryArtwork';
import type { RecoveryStateComponent } from './types/recoveryState';

/** Non-form recovery states, retry and accessible next actions. */
export const RecoveryState: RecoveryStateComponent = ({ view, purpose, invitation, slow, onRetry }) => {
  const checkingMessage = slow ? 'Still checking. This may take a moment…'
    : invitation ? 'Verifying your invitation securely…' : 'Verifying your reset link securely…';
  const showResetLink = !invitation && (view === 'missing' || view === 'expired');
  const showBack = view === 'missing' || view === 'expired' || view === 'attempts-exhausted';

  return (
    <div className="rdr-recovery__state">
      <RecoveryArtwork view={view} purpose={purpose}/>
      {view === 'checking' && <p className="rdr-recovery__slow" role="status">{checkingMessage}</p>}
      {view === 'success' && <Link className="rdr-login-v2__submit rdr-recovery__state-button" to="/login">Go to login</Link>}
      {view === 'unavailable' && <button className="rdr-login-v2__submit rdr-recovery__state-button"
        type="button" onClick={onRetry}>Try checking again</button>}
      {showResetLink && <Link className="rdr-login-v2__submit rdr-recovery__state-button"
        to="/forgot-password">Request new link</Link>}
      {showBack && <Link className="rdr-recovery__back" to="/login">
        <ArrowLeft size={16}/> Return to login
      </Link>}
    </div>
  );
};
