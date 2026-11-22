import { Box } from '@Components/components/Box';
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
    <Box as="div" className={`rdr-recovery__art rdr-recovery__art--${visualClass}`} role="status" aria-live="polite">
      <AuthStatusArt kind={statusKind}/>
      <Box as="strong" className="rdr-recovery__mobile-status-title">{mobileStatusTitle}</Box>
      <Box as="div" className="rdr-recovery__art-visual" aria-hidden="true">
        <Box as="span" className="rdr-recovery__art-orbit rdr-recovery__art-orbit--one"/>
        <Box as="span" className="rdr-recovery__art-orbit rdr-recovery__art-orbit--two"/>
        <Box as="span" className="rdr-recovery__art-icon"><Icon size={44} strokeWidth={1.55}/></Box>
        {success && <Sparkles className="rdr-recovery__sparkle" size={25}/>}
      </Box>
      <Box as="strong" className="rdr-recovery__art-title">{title}</Box>
      <Box as="p" className="rdr-recovery__art-message">{message}</Box>
    </Box>
  );
};
