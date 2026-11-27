import { Box, BOX_TAG } from '@Components/components/Box';

import { AuthStatusArt } from '../../../../components/authStatusArt';
import { getRecoveryArtworkCopy } from '../../helpers/recovery.getArtworkCopy';
import { getRecoveryArtworkDescriptor } from '../../helpers/recovery.getArtworkDescriptor';
import styles from './RecoveryArtwork.module.css';
import type { RecoveryArtworkComponent } from './types/recoveryArtwork';

/** One accessible status; the unused, permanently hidden legacy icon subtree is removed. */
export const RecoveryArtwork: RecoveryArtworkComponent = ({ view, purpose }) => {
  const { statusKind } = getRecoveryArtworkDescriptor(view);
  const { title, mobileStatusTitle, message } = getRecoveryArtworkCopy({ state: view, purpose });

  return (
    <Box as={BOX_TAG.DIV} className={styles.root} role="status" aria-live="polite">
      <AuthStatusArt kind={statusKind} />
      <Box as={BOX_TAG.STRONG} className={styles.mobileStatusTitle}>
        {mobileStatusTitle}
      </Box>
      <Box as={BOX_TAG.STRONG} className={styles.title}>
        {title}
      </Box>
      <Box as={BOX_TAG.P} className={styles.message}>
        {message}
      </Box>
    </Box>
  );
};
