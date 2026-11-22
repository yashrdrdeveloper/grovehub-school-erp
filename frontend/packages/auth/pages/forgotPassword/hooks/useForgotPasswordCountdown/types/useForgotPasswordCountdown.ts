export interface ForgotPasswordCountdownModel {
  readonly secondsLeft: number;
  readonly startCountdown: () => void;
  readonly isCoolingDown: () => boolean;
}

export type UseForgotPasswordCountdown = () => ForgotPasswordCountdownModel;
