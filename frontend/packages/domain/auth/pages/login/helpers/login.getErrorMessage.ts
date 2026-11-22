import { _find } from '@Base/utilities/find';
import { LOGIN_COPY } from '../constants/login.copy';
import { LOGIN_AUTH_ERROR_RULES } from '../constants/login.errorRules';
import type { GetLoginErrorMessage } from './types/login.getErrorMessage';

/** Never expose provider messages, technical details, or credentials. */
export const getLoginErrorMessage: GetLoginErrorMessage = (failure) => {
  if (!(failure instanceof Error)) return LOGIN_COPY.unavailable;
  const match = _find(LOGIN_AUTH_ERROR_RULES, rule => rule.pattern.test(failure.message));
  if (match) return match.message;
  return LOGIN_COPY.unavailable;
};
