import type { ReauthenticateRequest } from './reauthenticate.request';
import type { StatusResult } from './reauthenticate.response';

export type Reauthenticate = (password: ReauthenticateRequest['password']) => Promise<StatusResult>;
