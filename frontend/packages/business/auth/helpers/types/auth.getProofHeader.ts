import type { AUTH_PROOF_SCHEME } from '../../constants/auth.proofScheme';

export type AuthProofPurpose = keyof typeof AUTH_PROOF_SCHEME;
export type GetAuthProofHeader = (purpose: AuthProofPurpose, proof: string) => {
  readonly Authorization: string;
};
