import type { ExchangeRecoveryRequest } from './exchangeRecovery.request';
import type { GrantResult } from './exchangeRecovery.response';

export type ExchangeRecovery = (proof: ExchangeRecoveryRequest['proof']) => Promise<GrantResult>;
