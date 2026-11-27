import type { ValidateRecoveryRequest } from './validateRecovery.request';
import type { GrantResult } from './validateRecovery.response';

export type ValidateRecovery = (proof: ValidateRecoveryRequest['proof']) => Promise<GrantResult>;
