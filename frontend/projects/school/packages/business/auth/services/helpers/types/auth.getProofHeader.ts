import type { AuthProofPurpose } from '../../constants/auth.proofScheme';

export type GetAuthProofHeader = (
  purpose: AuthProofPurpose,
  proof: string,
) => {
  readonly Authorization: string;
};
