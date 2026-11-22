import { apiRequest } from '@Business/api';

import type { Login } from './types/login';
import type { LoginResponse } from './types/login.response';

const LOGIN_ENDPOINT = '/identity/public/auth/login';

/** Identity login transport. Session cookies remain owned by the browser/API. */
export const login: Login = (identifier, password) =>
  apiRequest<LoginResponse>(LOGIN_ENDPOINT, {
    method: 'POST',
    body: { identifier, password },
  });
