export interface CreateAccountResponse {
  readonly accountId: string;
  readonly personGid: string;
  readonly email: string;
  readonly status: string;
  readonly invitation: { readonly status: string };
}

export interface AccountSetupStatus {
  readonly verified: boolean;
  readonly expiresIn: number;
}

export interface ResendInvitationResponse {
  readonly accountId: string;
  readonly invitationStatus: string;
}
