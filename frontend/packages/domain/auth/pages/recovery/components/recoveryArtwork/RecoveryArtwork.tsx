import { AlertCircle, Clock3, Link2Off, PartyPopper, ShieldCheck, Sparkles } from '@Components/icons';

import { AuthStatusArt } from '../../../../components/authStatusArt';
import { getRecoveryArtworkCopy } from '../../helpers/recovery.getArtworkCopy';
import type { RecoveryArtworkComponent } from './types/recoveryArtwork';

/** State rendering and wording are separate from recovery operations. */
export const RecoveryArtwork: RecoveryArtworkComponent = ({ view, purpose }) => {
  const success = view === 'success';
  const checking = view === 'checking';
  const Icon = success ? PartyPopper : checking ? ShieldCheck
    : view === 'missing' ? Link2Off : view === 'attempts-exhausted' ? AlertCircle : Clock3;
  const statusKind = success ? 'success' : checking ? 'checking'
    : view === 'expired' ? 'expired' : view === 'unavailable' ? 'unavailable'
      : view === 'attempts-exhausted' ? 'attempts-exhausted' : 'missing';
  const { title, mobileStatusTitle, message } = getRecoveryArtworkCopy({ state: view, purpose });
  const visualClass = success ? 'success' : checking ? 'checking' : 'problem';

  return (
    <div className={`rdr-recovery__art rdr-recovery__art--${visualClass}`} role="status" aria-live="polite">
      <AuthStatusArt kind={statusKind}/>
      <strong className="rdr-recovery__mobile-status-title">{mobileStatusTitle}</strong>
      <div className="rdr-recovery__art-visual" aria-hidden="true">
        <span className="rdr-recovery__art-orbit rdr-recovery__art-orbit--one"/>
        <span className="rdr-recovery__art-orbit rdr-recovery__art-orbit--two"/>
        <span className="rdr-recovery__art-icon"><Icon size={44} strokeWidth={1.55}/></span>
        {success && <Sparkles className="rdr-recovery__sparkle" size={25}/>}
      </div>
      <strong className="rdr-recovery__art-title">{title}</strong>
      <p className="rdr-recovery__art-message">{message}</p>
    </div>
  );
};
