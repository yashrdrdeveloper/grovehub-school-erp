import { LOGIN_AUTH_ERROR_RULES, LOGIN_COPY } from '../config';

/** Translate provider failures into stable, user-safe form messages. */
export const getLoginErrorMessage = (failure: unknown): string => {
  const message = failure instanceof Error ? failure.message : '';
  return LOGIN_AUTH_ERROR_RULES.find(rule => rule.pattern.test(message))?.message
    ?? LOGIN_COPY.unavailable;
};
