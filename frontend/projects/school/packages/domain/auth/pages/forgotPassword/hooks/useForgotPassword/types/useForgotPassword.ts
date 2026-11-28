import type { FormEventHandler } from 'react';
import type { ForgotPasswordState, ForgotPasswordViewModel } from '../../../types';

export interface ForgotPasswordController {
  readonly state: ForgotPasswordState;
  readonly view: ForgotPasswordViewModel;
  readonly onSubmit: FormEventHandler<HTMLFormElement>;
  readonly onEmailChange: (email: ForgotPasswordState['email']) => void;
  readonly onUseAnotherEmail: () => void;
}

export type UseForgotPassword = () => ForgotPasswordController;
