import { Eye, EyeOff } from '@Components/icons';
import type { GetLoginPasswordViewModel } from './types/login.getPasswordViewModel';

export const getLoginPasswordViewModel: GetLoginPasswordViewModel = (visible) => {
  if (visible) return { Icon: EyeOff, label: 'Hide password', type: 'text' };
  return { Icon: Eye, label: 'Show password', type: 'password' };
};
