import { LoginField } from './login.fields';

/** Pure data: executable validation belongs in helpers. */
export enum LoginValidationKind {
  REQUIRED = 'REQUIRED',
  EMAIL_FORMAT = 'EMAIL_FORMAT',
}

export const LOGIN_EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

export const LOGIN_VALIDATION_RULES: Readonly<Record<LoginField, readonly LoginValidationKind[]>> =
  Object.freeze({
    [LoginField.EMAIL]: Object.freeze([
      LoginValidationKind.REQUIRED,
      LoginValidationKind.EMAIL_FORMAT,
    ]),
    [LoginField.PASSWORD]: Object.freeze([LoginValidationKind.REQUIRED]),
  });

export const LOGIN_REQUIRED_MESSAGES: Readonly<Record<LoginField, string>> = Object.freeze({
  [LoginField.EMAIL]: 'Please enter your email address.',
  [LoginField.PASSWORD]: 'Please enter your password.',
});

export const LOGIN_INVALID_EMAIL_MESSAGE = 'Enter a valid email address.';
