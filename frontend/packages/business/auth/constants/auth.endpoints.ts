import { HTTP_METHOD } from '../../api/constants/api.method';

import { AuthOperation } from './auth.operations';

import type { HttpMethod } from '../../api/constants/types/api.method';

export const AUTH_API_PREFIX = Object.freeze({
  PUBLIC: '/identity/public/auth',
  PROTECTED: '/identity/protected/auth',
} as const);

export interface AuthEndpoint {
  readonly method: HttpMethod;
  readonly path: string;
}

/** Every Identity endpoint is defined once and indexed by its operation enum. */
export const AUTH_ENDPOINT = {
  [AuthOperation.LOGIN]: {
    method: HTTP_METHOD.POST,
    path: `${AUTH_API_PREFIX.PUBLIC}/login`,
  },
  [AuthOperation.REFRESH]: {
    method: HTTP_METHOD.POST,
    path: `${AUTH_API_PREFIX.PUBLIC}/refresh`,
  },
  [AuthOperation.FORGOT_PASSWORD]: {
    method: HTTP_METHOD.POST,
    path: `${AUTH_API_PREFIX.PUBLIC}/forgot-password`,
  },
  [AuthOperation.VALIDATE_RECOVERY]: {
    method: HTTP_METHOD.GET,
    path: `${AUTH_API_PREFIX.PUBLIC}/recovery/validate`,
  },
  [AuthOperation.VALIDATE_INVITATION]: {
    method: HTTP_METHOD.GET,
    path: `${AUTH_API_PREFIX.PUBLIC}/invitations/validate`,
  },
  [AuthOperation.EXCHANGE_RECOVERY]: {
    method: HTTP_METHOD.POST,
    path: `${AUTH_API_PREFIX.PUBLIC}/recovery/exchange`,
  },
  [AuthOperation.EXCHANGE_INVITATION]: {
    method: HTTP_METHOD.POST,
    path: `${AUTH_API_PREFIX.PUBLIC}/invitations/exchange`,
  },
  [AuthOperation.RESET_GRANT_STATUS]: {
    method: HTTP_METHOD.GET,
    path: `${AUTH_API_PREFIX.PUBLIC}/reset-password/status`,
  },
  [AuthOperation.RESET_PASSWORD]: {
    method: HTTP_METHOD.POST,
    path: `${AUTH_API_PREFIX.PUBLIC}/reset-password`,
  },
  [AuthOperation.ME]: {
    method: HTTP_METHOD.GET,
    path: `${AUTH_API_PREFIX.PROTECTED}/me`,
  },
  [AuthOperation.LOGOUT]: {
    method: HTTP_METHOD.POST,
    path: `${AUTH_API_PREFIX.PROTECTED}/logout`,
  },
  [AuthOperation.LOGOUT_ALL]: {
    method: HTTP_METHOD.POST,
    path: `${AUTH_API_PREFIX.PROTECTED}/logout-all`,
  },
  [AuthOperation.CHANGE_PASSWORD]: {
    method: HTTP_METHOD.POST,
    path: `${AUTH_API_PREFIX.PROTECTED}/change-password`,
  },
  [AuthOperation.REAUTHENTICATE]: {
    method: HTTP_METHOD.POST,
    path: `${AUTH_API_PREFIX.PROTECTED}/reauthenticate`,
  },
  [AuthOperation.SESSIONS]: {
    method: HTTP_METHOD.GET,
    path: `${AUTH_API_PREFIX.PROTECTED}/sessions`,
  },
  [AuthOperation.REVOKE_SESSION]: {
    method: HTTP_METHOD.DELETE,
    path: `${AUTH_API_PREFIX.PROTECTED}/sessions/`,
  },
} as const satisfies Record<AuthOperation, AuthEndpoint>;
