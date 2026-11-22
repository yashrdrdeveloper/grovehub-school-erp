import type { Dispatch, SetStateAction } from 'react';
import type { RecoveryFieldErrors } from '../../../../../../types';

export interface RecoveryFormModel {
  readonly password: string;
  readonly confirmation: string;
  readonly visible: boolean;
  readonly error: string;
  readonly fieldErrors: RecoveryFieldErrors;
  readonly setError: Dispatch<SetStateAction<string>>;
  readonly setFieldErrors: Dispatch<SetStateAction<RecoveryFieldErrors>>;
  readonly clearForm: () => void;
  readonly clearCredentials: () => void;
  readonly onPasswordChange: (value: string) => void;
  readonly onConfirmationChange: (value: string) => void;
  readonly onToggleVisibility: () => void;
}
export type UseRecoveryForm = () => RecoveryFormModel;
