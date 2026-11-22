import type { FormEventHandler } from 'react';

import type { LoginFormModel } from '../../useLoginForm';

export interface LoginSubmissionModel {
  readonly isSubmitting: boolean;
  readonly onSubmit: FormEventHandler<HTMLFormElement>;
}

export type UseLoginSubmission = (model: LoginFormModel) => LoginSubmissionModel;
