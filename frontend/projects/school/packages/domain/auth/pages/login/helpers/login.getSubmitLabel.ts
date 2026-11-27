import { LOGIN_COPY } from '../constants/login.copy';
import type { GetLoginSubmitLabel } from './types/login.getSubmitLabel';

export const getLoginSubmitLabel: GetLoginSubmitLabel = (submitting) => {
  if (submitting) return LOGIN_COPY.busyButton;
  return LOGIN_COPY.button;
};
