import type { ForgotPasswordState } from '../../types';
export type GetForgotPasswordNoticeCopy = (
  state: ForgotPasswordState,
  emailChanged: boolean,
  noticeSuccess: boolean,
) => { readonly noticeTitle: string; readonly noticeMessage: string };
