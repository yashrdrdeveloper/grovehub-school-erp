/** Only these proof purposes are valid for Identity's bearer-grant exchange. */
export enum AuthProofPurpose {
  RECOVERY = 'RECOVERY',
  INVITATION = 'INVITATION',
}

export const AUTH_PROOF_SCHEME: Readonly<Record<AuthProofPurpose, string>> = {
  [AuthProofPurpose.RECOVERY]: 'Recovery',
  [AuthProofPurpose.INVITATION]: 'Invitation',
};
