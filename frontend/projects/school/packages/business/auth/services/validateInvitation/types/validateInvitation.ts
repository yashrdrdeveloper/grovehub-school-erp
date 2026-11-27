import type { ValidateInvitationRequest } from './validateInvitation.request';
import type { GrantResult } from './validateInvitation.response';

export type ValidateInvitation = (
  proof: ValidateInvitationRequest['proof'],
) => Promise<GrantResult>;
