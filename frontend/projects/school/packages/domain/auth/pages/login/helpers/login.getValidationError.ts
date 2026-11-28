import { LoginField } from '../constants/login.fields';
import {
  LOGIN_INVALID_EMAIL_MESSAGE,
  LOGIN_REQUIRED_MESSAGES,
  LOGIN_VALIDATION_RULES,
  LoginValidationKind,
} from '../constants/login.validation';
import { isValidLoginEmail } from './login.isValidEmail';
import type { GetLoginValidationError } from './types/login.getValidationError';

/** Evaluate a declarative rule list without executable functions in constants. */
export const getLoginValidationError: GetLoginValidationError = (field, value) => {
  for (const rule of LOGIN_VALIDATION_RULES[field]) {
    if (rule === LoginValidationKind.REQUIRED && value.length === 0) {
      return LOGIN_REQUIRED_MESSAGES[field];
    }
    if (rule === LoginValidationKind.EMAIL_FORMAT && field === LoginField.EMAIL) {
      if (!isValidLoginEmail(value)) return LOGIN_INVALID_EMAIL_MESSAGE;
    }
  }
  return undefined;
};
