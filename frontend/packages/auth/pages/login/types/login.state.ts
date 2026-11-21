/** Login-only view state. No server DTOs or provider types leak into the UI. */
export type LoginField = 'email' | 'password';

export type LoginFieldErrors = Partial<Record<LoginField, string>>;

export interface LoginState {
  readonly email: string;
  readonly password: string;
  readonly isPasswordVisible: boolean;
  readonly isSubmitting: boolean;
  readonly errorMessage: string;
  readonly fieldErrors: LoginFieldErrors;
}

export type LoginCredentials = Pick<LoginState, 'email' | 'password'>;
