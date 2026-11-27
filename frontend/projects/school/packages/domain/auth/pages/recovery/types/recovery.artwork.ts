import type { RecoveryPurpose, RecoveryViewState } from './recovery.props';

export interface RecoveryArtworkCopy {
  readonly mobileStatusTitle: string;
  readonly title: string;
  readonly message: string;
}

export interface RecoveryArtworkParams {
  readonly state: RecoveryViewState;
  readonly purpose: RecoveryPurpose;
}
