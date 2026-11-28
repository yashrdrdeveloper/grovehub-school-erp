import { LoginField } from './login.fields';

export const LOGIN_FIELD_ORDER = Object.freeze([
  LoginField.EMAIL,
  LoginField.PASSWORD,
] as const satisfies readonly LoginField[]);
