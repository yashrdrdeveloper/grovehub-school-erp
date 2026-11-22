import { RECOVERY_ARTWORK } from '../constants/recovery.artwork';
import type { GetRecoveryArtworkDescriptor } from './types/recovery.getArtworkDescriptor';

export const getRecoveryArtworkDescriptor: GetRecoveryArtworkDescriptor = (state) =>
  RECOVERY_ARTWORK[state];
