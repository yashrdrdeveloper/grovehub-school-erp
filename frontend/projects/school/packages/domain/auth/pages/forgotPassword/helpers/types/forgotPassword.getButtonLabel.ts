import type { ForgotPasswordState } from '../../types';
export type GetForgotPasswordButtonLabel = (
  state: ForgotPasswordState,
  noticeSuccess: boolean,
) => string;
