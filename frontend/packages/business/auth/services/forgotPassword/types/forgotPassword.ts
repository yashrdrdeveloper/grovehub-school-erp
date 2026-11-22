import type { ForgotPasswordRequest } from './forgotPassword.request';
import type { ForgotPasswordResponse } from './forgotPassword.response';

export type ForgotPassword = (email: ForgotPasswordRequest['email']) => Promise<ForgotPasswordResponse>;
