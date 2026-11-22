import type { FormEventHandler, ReactNode } from 'react';

import type { LoginState } from './login.state';

/** Data and callbacks supplied by the login container to either template. */
export interface LoginProps {
  readonly fieldId: string;
  readonly state: LoginState;
  readonly onEmailChange: (email: string) => void;
  readonly onPasswordChange: (password: string) => void;
  readonly onTogglePassword: () => void;
  readonly onSubmit: FormEventHandler<HTMLFormElement>;
}

/** Optional slots allow desktop/mobile composition without duplicating controls. */
export interface LoginFormProps extends LoginProps {
  readonly passwordLabelAction?: ReactNode;
  readonly betweenFields?: ReactNode;
}

/** Pure page composition props supplied by LoginContainer. */
export interface LoginPageProps {
  readonly controls: LoginProps;
  readonly isMobile: boolean;
}
