import type { FormEventHandler, ReactElement } from 'react';

export type RecoveryPurpose = 'recovery' | 'invitation';
export type RecoveryViewState =
  | 'missing' | 'checking' | 'expired' | 'unavailable'
  | 'ready' | 'saving' | 'success' | 'attempts-exhausted';

export interface RecoveryFieldErrors {
  password?: string;
  confirmation?: string;
}

export interface RecoveryPageProps {
  readonly purpose: RecoveryPurpose;
}

export interface RecoveryProps extends RecoveryPageProps {
  readonly invitation: boolean;
  readonly id: string;
  readonly view: RecoveryViewState;
  readonly slow: boolean;
  readonly password: string;
  readonly confirmation: string;
  readonly visible: boolean;
  readonly error: string;
  readonly fieldErrors: RecoveryFieldErrors;
  readonly attempts: number;
  readonly onPasswordChange: (value: string) => void;
  readonly onConfirmationChange: (value: string) => void;
  readonly onToggleVisibility: () => void;
  readonly onSubmit: FormEventHandler<HTMLFormElement>;
  readonly onRetry: () => void;
}

export type RecoveryComponent = (props: RecoveryProps) => ReactElement;
