import type { Dispatch, SetStateAction } from 'react';

import type { LoginFieldErrors, LoginState } from '../../../types';

/** All state mutators are local to the login page and never cross business boundaries. */
export interface LoginFormModel {
  readonly email: LoginState['email'];
  readonly password: LoginState['password'];
  readonly isPasswordVisible: LoginState['isPasswordVisible'];
  readonly isSubmitting: LoginState['isSubmitting'];
  readonly errorMessage: LoginState['errorMessage'];
  readonly fieldErrors: LoginFieldErrors;
  readonly setEmail: Dispatch<SetStateAction<LoginState['email']>>;
  readonly setPassword: Dispatch<SetStateAction<LoginState['password']>>;
  readonly setPasswordVisible: Dispatch<SetStateAction<LoginState['isPasswordVisible']>>;
  readonly setSubmitting: Dispatch<SetStateAction<LoginState['isSubmitting']>>;
  readonly setErrorMessage: Dispatch<SetStateAction<LoginState['errorMessage']>>;
  readonly setFieldErrors: Dispatch<SetStateAction<LoginFieldErrors>>;
}

export type UseLoginForm = () => LoginFormModel;
