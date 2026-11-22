import { LOGIN_AUTH_ERROR_RULES, LOGIN_COPY } from '../../config';

import type { GetLoginErrorMessage } from './types/getLoginErrorMessage';

/** Translate provider failures into stable, user-safe form messages. */
export const getLoginErrorMessage: GetLoginErrorMessage = (failure) => {
  const message = failure instanceof Error ? failure.message : '';
  return LOGIN_AUTH_ERROR_RULES.find(rule => rule.pattern.test(message))?.message
    ?? LOGIN_COPY.unavailable;
};
