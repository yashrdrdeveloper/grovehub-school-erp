import { Eye, EyeOff } from '@Components/icons';

/** Pure presentation selection; both labels are literal translator messages. */
export const getPasswordFieldVisibility = (visible: boolean) => {
  if (visible) {
    return {
      VisibilityIcon: EyeOff,
      toggleLabel: __('Hide password'),
      passwordType: 'text' as const,
    };
  }
  return {
    VisibilityIcon: Eye,
    toggleLabel: __('Show password'),
    passwordType: 'password' as const,
  };
};
