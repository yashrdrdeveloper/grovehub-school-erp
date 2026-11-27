/** State-to-copy projection; never mutate a label in the JSX component. */
export const getRecoveryPasswordSubmitLabel = (saving: boolean): string => {
  if (saving) return __('Securing your account…');
  return __('Update password');
};
