import { AUTH_ENDPOINT } from '../constants/auth.endpoints';

import type { GetAuthEndpoint } from './types/auth.getEndpoint';

export const getAuthEndpoint: GetAuthEndpoint = (operation) =>
  AUTH_ENDPOINT[operation];
