import { AuthOperation } from './constants/auth.operations';
import { AUTHENTICATION_OPERATIONS } from './auth.operations';

/** Preserve the public named API while dispatching via typed operation identifiers. */
export const authenticationApi = {
  login: AUTHENTICATION_OPERATIONS[AuthOperation.LOGIN],
  refresh: AUTHENTICATION_OPERATIONS[AuthOperation.REFRESH],
  forgotPassword: AUTHENTICATION_OPERATIONS[AuthOperation.FORGOT_PASSWORD],
  validateRecovery: AUTHENTICATION_OPERATIONS[AuthOperation.VALIDATE_RECOVERY],
  validateInvitation: AUTHENTICATION_OPERATIONS[AuthOperation.VALIDATE_INVITATION],
  exchangeRecovery: AUTHENTICATION_OPERATIONS[AuthOperation.EXCHANGE_RECOVERY],
  exchangeInvitation: AUTHENTICATION_OPERATIONS[AuthOperation.EXCHANGE_INVITATION],
  resetGrantStatus: AUTHENTICATION_OPERATIONS[AuthOperation.RESET_GRANT_STATUS],
  resetPassword: AUTHENTICATION_OPERATIONS[AuthOperation.RESET_PASSWORD],
  me: AUTHENTICATION_OPERATIONS[AuthOperation.ME],
  logout: AUTHENTICATION_OPERATIONS[AuthOperation.LOGOUT],
  logoutAll: AUTHENTICATION_OPERATIONS[AuthOperation.LOGOUT_ALL],
  changePassword: AUTHENTICATION_OPERATIONS[AuthOperation.CHANGE_PASSWORD],
  reauthenticate: AUTHENTICATION_OPERATIONS[AuthOperation.REAUTHENTICATE],
  sessions: AUTHENTICATION_OPERATIONS[AuthOperation.SESSIONS],
  revokeSession: AUTHENTICATION_OPERATIONS[AuthOperation.REVOKE_SESSION],
} as const;
