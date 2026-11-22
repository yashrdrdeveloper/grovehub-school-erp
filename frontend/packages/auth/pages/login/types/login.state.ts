import type { AccountEntity } from '@Business/auth';

/** Login-only view state. Reuse authoritative business fields, not full DTOs. */
export type LoginField = 'email' | 'password';

export type LoginFieldErrors = Partial<Record<LoginField, string>>;

export interface LoginState {
  readonly email: AccountEntity['email'];
  readonly password: string;
  readonly isPasswordVisible: boolean;
  readonly isSubmitting: boolean;
  readonly errorMessage: string;
  readonly fieldErrors: LoginFieldErrors;
}

export type LoginCredentials = Pick<LoginState, 'email' | 'password'>;
