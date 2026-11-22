import { apiRequest } from '@Business/api';

import type { Login } from './types/login';

const LOGIN_ENDPOINT = '/identity/public/auth/login';

/** Identity login transport. Session cookies remain owned by the browser/API. */
export const login: Login = (identifier, password) =>
  apiRequest(LOGIN_ENDPOINT, {
    method: 'POST',
    body: { identifier, password },
  });
