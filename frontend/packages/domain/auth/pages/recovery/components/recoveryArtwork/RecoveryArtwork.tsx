import { Box, BOX_TAG } from '@Components/components/Box';
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
    <Box
      as={BOX_TAG.DIV}
      className={`rdr-recovery__art rdr-recovery__art--${visualClass}`}
      role="status"
      aria-live="polite"
    >
      <AuthStatusArt kind={statusKind} />
      <Box as={BOX_TAG.STRONG} className="rdr-recovery__mobile-status-title">
        {mobileStatusTitle}
      </Box>
      <Box as={BOX_TAG.DIV} className="rdr-recovery__art-visual" aria-hidden="true">
        <Box as={BOX_TAG.SPAN} className="rdr-recovery__art-orbit rdr-recovery__art-orbit--one" />
        <Box as={BOX_TAG.SPAN} className="rdr-recovery__art-orbit rdr-recovery__art-orbit--two" />
        <Box as={BOX_TAG.SPAN} className="rdr-recovery__art-icon">
          <Icon size={44} strokeWidth={1.55} />
        </Box>
        {success && <Sparkles className="rdr-recovery__sparkle" size={25} />}
      </Box>
      <Box as={BOX_TAG.STRONG} className="rdr-recovery__art-title">
        {title}
      </Box>
      <Box as={BOX_TAG.P} className="rdr-recovery__art-message">
        {message}
      </Box>
    </Box>
  );
};
