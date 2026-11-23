import { RECOVERY_ARTWORK_COPY } from '../constants/recovery.artworkCopy';
import type { GetRecoveryArtworkCopy } from './types/recovery.getArtworkCopy';

/** Pure purpose/state lookup; callers consume already translated messages. */
export const getRecoveryArtworkCopy: GetRecoveryArtworkCopy = ({ state, purpose }) =>
  RECOVERY_ARTWORK_COPY[purpose][state];
