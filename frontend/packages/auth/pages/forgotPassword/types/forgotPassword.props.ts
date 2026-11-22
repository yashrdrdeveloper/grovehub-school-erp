import type { FormEventHandler, ReactElement } from 'react';

import type { ForgotPasswordState } from './forgotPassword.state';
import type { ForgotPasswordViewModel } from './forgotPassword.viewModel';

export interface ForgotPasswordProps {
  readonly id: string;
  readonly formId: string;
  readonly state: ForgotPasswordState;
  readonly view: ForgotPasswordViewModel;
  readonly onSubmit: FormEventHandler<HTMLFormElement>;
  readonly onEmailChange: (email: ForgotPasswordState['email']) => void;
  readonly onUseAnotherEmail: () => void;
}

export type ForgotPasswordComponent = (props: ForgotPasswordProps) => ReactElement;
