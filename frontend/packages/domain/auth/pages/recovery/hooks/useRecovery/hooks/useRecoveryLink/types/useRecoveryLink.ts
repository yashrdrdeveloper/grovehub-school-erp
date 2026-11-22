import type { Dispatch, SetStateAction } from 'react';
import type { RecoveryPurpose, RecoveryViewState } from '../../../../../../types';

export interface RecoveryLinkModel {
  readonly view: RecoveryViewState;
  readonly slow: boolean;
  readonly setView: Dispatch<SetStateAction<RecoveryViewState>>;
  readonly onRetry: () => void;
}
export type UseRecoveryLink = (purpose: RecoveryPurpose, clearForm: () => void) => RecoveryLinkModel;
