export const getRecoveryCheckingMessage = (params: {
  readonly invitation: boolean;
  readonly slow: boolean;
}): string => {
  if (params.slow) return __('Still checking. This may take a moment…');
  if (params.invitation) return __('Verifying your invitation securely…');
  return __('Verifying your reset link securely…');
};
