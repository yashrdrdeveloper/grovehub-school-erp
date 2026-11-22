import { LOGIN_COPY } from './login.copy';

/** Declarative safe-message matching, no executable validation in constants. */
export const LOGIN_AUTH_ERROR_RULES = Object.freeze([
  Object.freeze({
    pattern: /invalid creden|incorrect|unauthoriz|invalid login|wrong password|\b401\b/i,
    message: LOGIN_COPY.incorrectCredentials,
  }),
]);
