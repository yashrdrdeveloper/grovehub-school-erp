import { LOGIN_COPY } from './login.copy';

/** Patterns are declarative mappings; safe message selection is a pure helper. */
export const LOGIN_AUTH_ERROR_RULES = Object.freeze([
  Object.freeze({
    pattern: /invalid creden|incorrect|unauthoriz|invalid login|wrong password|\\b401\\b/i,
    message: LOGIN_COPY.incorrectCredentials,
  }),
]);
