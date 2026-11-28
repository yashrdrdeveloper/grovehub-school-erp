import type { AccountEntity } from '@Business/auth';

export interface ForgotPasswordState {
  readonly email: AccountEntity['email'];
  readonly lastRequestedEmail: AccountEntity['email'];
  readonly busy: boolean;
  readonly secondsLeft: number;
  readonly error: string;
  readonly emailError: string;
}
