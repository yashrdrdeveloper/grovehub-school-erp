import type { FormEventHandler } from 'react';
import type { RecoveryPurpose } from '../../../../../types';
import type { RecoveryFormModel } from '../../useRecoveryForm/types/useRecoveryForm';
import type { RecoveryLinkModel } from '../../useRecoveryLink/types/useRecoveryLink';

export interface RecoverySubmissionModel {
  readonly onSubmit: FormEventHandler<HTMLFormElement>;
}
export type UseRecoverySubmission = (
  purpose: RecoveryPurpose,
  form: RecoveryFormModel,
  link: RecoveryLinkModel,
) => RecoverySubmissionModel;
