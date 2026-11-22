import type { ExchangeInvitationRequest } from './exchangeInvitation.request';
import type { GrantResult } from './exchangeInvitation.response';

export type ExchangeInvitation = (proof: ExchangeInvitationRequest['proof']) => Promise<GrantResult>;
