import { refresh } from './refresh';
import { forgotPassword } from './forgotPassword';
import { validateRecovery } from './validateRecovery';
import { validateInvitation } from './validateInvitation';
import { exchangeRecovery } from './exchangeRecovery';
import { exchangeInvitation } from './exchangeInvitation';
import { resetGrantStatus } from './resetGrantStatus';
import { resetPassword } from './resetPassword';
import { me } from './me';
import { logout } from './logout';
import { logoutAll } from './logoutAll';
import { changePassword } from './changePassword';
import { reauthenticate } from './reauthenticate';
import { sessions } from './sessions';
import { revokeSession } from './revokeSession';
import { login } from './login';

import { AuthOperation } from './constants/auth.operations';

/** Exhaustive operation-to-function registry; each handler owns its transport and contracts. */
export const AUTHENTICATION_OPERATIONS = {
  [AuthOperation.LOGIN]: login,
  [AuthOperation.REFRESH]: refresh,
  [AuthOperation.FORGOT_PASSWORD]: forgotPassword,
  [AuthOperation.VALIDATE_RECOVERY]: validateRecovery,
  [AuthOperation.VALIDATE_INVITATION]: validateInvitation,
  [AuthOperation.EXCHANGE_RECOVERY]: exchangeRecovery,
  [AuthOperation.EXCHANGE_INVITATION]: exchangeInvitation,
  [AuthOperation.RESET_GRANT_STATUS]: resetGrantStatus,
  [AuthOperation.RESET_PASSWORD]: resetPassword,
  [AuthOperation.ME]: me,
  [AuthOperation.LOGOUT]: logout,
  [AuthOperation.LOGOUT_ALL]: logoutAll,
  [AuthOperation.CHANGE_PASSWORD]: changePassword,
  [AuthOperation.REAUTHENTICATE]: reauthenticate,
  [AuthOperation.SESSIONS]: sessions,
  [AuthOperation.REVOKE_SESSION]: revokeSession,
} as const;
