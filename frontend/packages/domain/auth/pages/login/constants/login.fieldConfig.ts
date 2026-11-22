import { LoginField } from './login.fields';
import type { LoginFieldConfig } from './types/login.fieldConfig';

/** One immutable set of presentation fields keyed by the LoginField enum. */
export const LOGIN_FIELD_CONFIG: Readonly<Record<LoginField, Readonly<LoginFieldConfig>>> =
  Object.freeze({
    [LoginField.EMAIL]: Object.freeze({
      label: 'Email address',
      placeholder: 'Email',
      autoComplete: 'username',
    }),
    [LoginField.PASSWORD]: Object.freeze({
      label: 'Password',
      placeholder: 'Password',
      autoComplete: 'current-password',
    }),
  });
