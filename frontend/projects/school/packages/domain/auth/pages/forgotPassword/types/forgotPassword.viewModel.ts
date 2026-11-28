export interface ForgotPasswordViewModel {
  readonly cooldown: boolean;
  readonly noticeSuccess: boolean;
  readonly noticeTitle: string;
  readonly noticeMessage: string;
  readonly buttonLabel: string;
  readonly resendLabel: string;
}
