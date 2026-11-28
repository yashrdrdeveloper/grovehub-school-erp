import { Eye, EyeOff } from '@Components/icons';

/** Password and confirmation toggle share one visible state. */
export const getRecoveryPasswordVisibility = (visible: boolean) => {
  if (visible) {
    return {
      VisibilityIcon: EyeOff,
      passwordType: 'text' as const,
      visibilityLabel: __('Hide passwords'),
      confirmationLabel: __('Hide confirmation password'),
    };
  }
  return {
    VisibilityIcon: Eye,
    passwordType: 'password' as const,
    visibilityLabel: __('Show passwords'),
    confirmationLabel: __('Show confirmation password'),
  };
};
