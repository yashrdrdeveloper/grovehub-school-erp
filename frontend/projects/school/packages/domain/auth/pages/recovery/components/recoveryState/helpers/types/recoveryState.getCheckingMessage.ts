export type GetRecoveryCheckingMessage = (params: {
  readonly invitation: boolean;
  readonly slow: boolean;
}) => string;
