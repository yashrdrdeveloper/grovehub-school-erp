import type { LoginCredentials } from '../../../types';

export type NormalizeLoginEmail = (
  email: LoginCredentials['email']
) => LoginCredentials['email'];
