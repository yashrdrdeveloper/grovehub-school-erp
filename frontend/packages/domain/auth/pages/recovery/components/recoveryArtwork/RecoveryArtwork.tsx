import { Sparkles } from '@Components/icons';
import { getRecoveryArtworkDescriptor } from '../../helpers/recovery.getArtworkDescriptor';

import { AuthStatusArt } from '../../../../components/authStatusArt';
import { getRecoveryArtworkCopy } from '../../helpers/recovery.getArtworkCopy';
import type { RecoveryArtworkComponent } from './types/recoveryArtwork';

/** State rendering and wording are separate from recovery operations. */
export const RecoveryArtwork: RecoveryArtworkComponent = ({ view, purpose }) => {
  const success = view === 'success';
  const { Icon, statusKind, visualClass } = getRecoveryArtworkDescriptor(view);
  const { title, mobileStatusTitle, message } = getRecoveryArtworkCopy({ state: view, purpose });

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
