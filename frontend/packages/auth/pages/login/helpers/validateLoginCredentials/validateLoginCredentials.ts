import { LOGIN_FIELD_ORDER, LOGIN_VALIDATION_RULES } from '../../config';
import { normalizeLoginEmail } from '../normalizeLoginEmail';

import type { ValidateLoginCredentials } from './types/validateLoginCredentials';

/** Pure configuration-driven credential validation; never logs the password. */
export const validateLoginCredentials: ValidateLoginCredentials = (credentials) => {
  const normalizedEmail = normalizeLoginEmail(credentials.email);
  const values = { email: normalizedEmail, password: credentials.password };
  const fieldErrors: ReturnType<ValidateLoginCredentials>['fieldErrors'] = {};

  for (const field of LOGIN_FIELD_ORDER) {
    const failedRule = LOGIN_VALIDATION_RULES[field].find(rule => !rule.isValid(values[field]));
    if (failedRule) fieldErrors[field] = failedRule.message;
  }

  return { normalizedEmail, fieldErrors };
};
