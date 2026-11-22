import type { ValidateRecoveryPassword } from './types/validateRecoveryPassword';

/** Pure browser-side checks must not count toward server setup attempts. */
export const validateRecoveryPassword: ValidateRecoveryPassword = ({ password, confirmation }) => {
  const fieldErrors: ReturnType<ValidateRecoveryPassword> = {};
  if (!password) fieldErrors.password = 'Please enter a password.';
  else if (password.length < 12) fieldErrors.password = 'Use at least 12 characters.';
  if (!confirmation) fieldErrors.confirmation = 'Please confirm your password.';
  else if (confirmation !== password) fieldErrors.confirmation = 'Passwords do not match.';
  return fieldErrors;
};
