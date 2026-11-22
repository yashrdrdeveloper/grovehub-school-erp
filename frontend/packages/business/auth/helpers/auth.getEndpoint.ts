import { AUTH_ENDPOINT } from '../constants/auth.endpoints';

import type { AuthOperation } from '../constants/auth.operations';
import type { AuthEndpoint } from '../constants/auth.endpoints';

export const getAuthEndpoint = (operation: AuthOperation): AuthEndpoint =>
  AUTH_ENDPOINT[operation];
