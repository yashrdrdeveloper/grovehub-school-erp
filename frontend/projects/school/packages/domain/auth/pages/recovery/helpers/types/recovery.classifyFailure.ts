export type RecoveryFailureKind = 'attempts-exhausted' | 'expired' | 'retryable';
export type ClassifyRecoveryFailure = (failure: unknown) => RecoveryFailureKind;
