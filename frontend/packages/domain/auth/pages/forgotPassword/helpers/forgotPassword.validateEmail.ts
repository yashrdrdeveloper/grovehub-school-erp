import { _trim } from '@Base/utilities/trim';
import type { ValidateForgotPasswordEmail } from './types/forgotPassword.validateEmail';

const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

/** Pure input validation; neither existence nor eligibility is disclosed. */
export const validateForgotPasswordEmail: ValidateForgotPasswordEmail = (email) => {
  const normalizedEmail = _trim(email);
  if (!normalizedEmail) return { normalizedEmail, error: 'Please enter your email address.' };
  if (!EMAIL_PATTERN.test(normalizedEmail)) {
    return { normalizedEmail, error: 'Enter a valid email address.' };
  }
  return { normalizedEmail, error: '' };
};
