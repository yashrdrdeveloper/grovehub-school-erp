import { ApiError } from '@Business/api';
import type { ClassifyRecoveryFailure } from './types/recovery.classifyFailure';

export const classifyRecoveryFailure: ClassifyRecoveryFailure = failure => {
  if (!(failure instanceof ApiError)) return 'retryable';
  if (failure.code === 'IDENTITY_SETUP_ATTEMPTS_EXCEEDED') return 'attempts-exhausted';
  if (failure.code === 'IDENTITY_TOKEN_INVALID' || [401, 403, 410].includes(failure.status)) {
    return 'expired';
  }
  return 'retryable';
};
