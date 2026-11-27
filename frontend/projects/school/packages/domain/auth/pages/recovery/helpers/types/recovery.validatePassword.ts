import type { RecoveryFieldErrors, RecoveryProps } from '../../types';

export type RecoveryCredentials = Pick<RecoveryProps, 'password' | 'confirmation'>;
export type ValidateRecoveryPassword = (credentials: RecoveryCredentials) => RecoveryFieldErrors;
