import type { LoginField } from '../types';

type LoginValidationRule = Readonly<{
  isValid: (value: string) => boolean;
  message: string;
}>;

export const LOGIN_COPY = Object.freeze({
  desktopTitle: 'Welcome back.',
  desktopDescription: 'Sign in to continue to your RDR School workspace.',
  mobileTitle: 'Welcome back.',
  mobileDescription: 'Sign in to your RDR School workspace.',
  button: 'Sign in',
  busyButton: 'Signing in…',
  forgotPassword: 'Forgot password?',
  incorrectCredentials: 'Incorrect email or password. Please try again.',
  unavailable: 'Unable to sign in right now. Please try again.',
});

export const LOGIN_FIELDS = Object.freeze({
  email: Object.freeze({
    label: 'Email address',
    placeholder: 'Email',
    autoComplete: 'username',
  }),
  password: Object.freeze({
    label: 'Password',
    placeholder: 'Password',
    autoComplete: 'current-password',
  }),
});

export const LOGIN_FIELD_ORDER = Object.freeze(['email', 'password'] as const satisfies readonly LoginField[]);

const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

export const LOGIN_VALIDATION_RULES: Readonly<Record<LoginField, readonly LoginValidationRule[]>> = {
  email: [
    { isValid: value => value.length > 0, message: 'Please enter your email address.' },
    { isValid: value => EMAIL_PATTERN.test(value), message: 'Enter a valid email address.' },
  ],
  password: [
    { isValid: value => value.length > 0, message: 'Please enter your password.' },
  ],
};

export const LOGIN_AUTH_ERROR_RULES = Object.freeze([
  { pattern: /invalid creden|incorrect|unauthoriz|invalid login|wrong password|\b401\b/i, message: LOGIN_COPY.incorrectCredentials },
]);

export const LOGIN_ROUTES = Object.freeze({
  success: '/dashboard',
  forgotPassword: '/forgot-password',
});
