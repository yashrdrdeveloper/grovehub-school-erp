import type { LoginCredentials, LoginFieldErrors } from '../../types';

export interface ValidateLoginCredentialsResult {
  readonly normalizedEmail: LoginCredentials['email'];
  readonly fieldErrors: LoginFieldErrors;
}

export type ValidateLoginCredentials = (
  credentials: LoginCredentials,
) => ValidateLoginCredentialsResult;
