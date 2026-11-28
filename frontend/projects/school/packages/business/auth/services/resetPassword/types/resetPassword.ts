import type { ResetPasswordRequest } from './resetPassword.request';
import type { StatusResult } from './resetPassword.response';

export type ResetPassword = (
  newPassword: ResetPasswordRequest['newPassword'],
) => Promise<StatusResult>;
