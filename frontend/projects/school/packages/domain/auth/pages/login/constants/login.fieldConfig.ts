import { LoginField } from './login.fields';
import type { LoginFieldConfig } from './types/login.fieldConfig';

/** One immutable set of presentation fields keyed by the LoginField enum. */
export const LOGIN_FIELD_CONFIG: Readonly<Record<LoginField, Readonly<LoginFieldConfig>>> =
  Object.freeze({
    [LoginField.EMAIL]: Object.freeze({
      label: __('Email address'),
      placeholder: __('Email'),
      autoComplete: 'username',
    }),
    [LoginField.PASSWORD]: Object.freeze({
      label: __('Password'),
      placeholder: __('Password'),
      autoComplete: 'current-password',
    }),
  });
