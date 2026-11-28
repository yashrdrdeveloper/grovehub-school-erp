import type { LoginRequest } from './login.request';
import type { LoginResponse } from './login.response';

export type Login = (
  identifier: LoginRequest['identifier'],
  password: LoginRequest['password'],
) => Promise<LoginResponse>;
