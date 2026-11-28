import type { ChangePasswordRequest } from './changePassword.request';
import type { StatusResult } from './changePassword.response';

export type ChangePassword = (
  newPassword: ChangePasswordRequest['newPassword'],
) => Promise<StatusResult>;
