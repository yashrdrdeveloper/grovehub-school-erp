import { AUTH_PROOF_SCHEME } from '../constants/auth.proofScheme';

import type { GetAuthProofHeader } from './types/auth.getProofHeader';

export const getAuthProofHeader: GetAuthProofHeader = (purpose, proof) => ({
  Authorization: `${AUTH_PROOF_SCHEME[purpose]} ${proof}`,
});
