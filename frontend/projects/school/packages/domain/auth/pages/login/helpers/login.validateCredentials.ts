import { _trim } from '@Base/utilities/trim';
import { LoginField } from '../constants/login.fields';
import { LOGIN_FIELD_ORDER } from '../constants/login.fieldOrder';
import { getLoginValidationError } from './login.getValidationError';
import type { ValidateLoginCredentials } from './types/login.validateCredentials';

/** Trim only the identifier; password bytes pass through unchanged. */
export const validateLoginCredentials: ValidateLoginCredentials = (credentials) => {
  const normalizedEmail = _trim(credentials.email);
  const values: Record<LoginField, string> = {
    [LoginField.EMAIL]: normalizedEmail,
    [LoginField.PASSWORD]: credentials.password,
  };
  const fieldErrors: ReturnType<ValidateLoginCredentials>['fieldErrors'] = {};

  for (const field of LOGIN_FIELD_ORDER) {
    const message = getLoginValidationError(field, values[field]);
    if (message) fieldErrors[field] = message;
  }

  return { normalizedEmail, fieldErrors };
};
