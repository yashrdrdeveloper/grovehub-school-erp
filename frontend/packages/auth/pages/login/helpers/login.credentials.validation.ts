import { LOGIN_FIELD_ORDER, LOGIN_VALIDATION_RULES } from '../config';
import type { LoginCredentials, LoginFieldErrors } from '../types';

/** Pure, configuration-driven credential validation; never logs the password. */
export const validateLoginCredentials = (credentials: LoginCredentials): {
  normalizedEmail: string;
  fieldErrors: LoginFieldErrors;
} => {
  const normalizedEmail = credentials.email.trim();
  const values = { email: normalizedEmail, password: credentials.password };
  const fieldErrors: LoginFieldErrors = {};

  for (const field of LOGIN_FIELD_ORDER) {
    const failedRule = LOGIN_VALIDATION_RULES[field].find(rule => !rule.isValid(values[field]));
    if (failedRule) fieldErrors[field] = failedRule.message;
  }

  return { normalizedEmail, fieldErrors };
};
