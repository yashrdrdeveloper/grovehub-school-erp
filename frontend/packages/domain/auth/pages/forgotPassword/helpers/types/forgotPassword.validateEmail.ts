import type { ForgotPasswordState } from '../../types';

export interface ValidateForgotPasswordEmailResult {
  readonly normalizedEmail: ForgotPasswordState['email'];
  readonly error: ForgotPasswordState['emailError'];
}

export type ValidateForgotPasswordEmail = (
  email: ForgotPasswordState['email']
) => ValidateForgotPasswordEmailResult;
