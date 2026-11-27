import type { LoginField } from '../../constants/login.fields';

export type GetLoginValidationError = (field: LoginField, value: string) => string | undefined;
