import { DesktopLogin } from '../templates/desktopLogin';
import { MobileLogin } from '../templates/mobileLogin';
import type { GetLoginFormTemplate } from './types/login.formTemplate';

/** Pure choice of one mounted form; no mutable JSX component selection in Login.tsx. */
export const getLoginFormTemplate: GetLoginFormTemplate = (isMobile) => {
  if (isMobile) return MobileLogin;
  return DesktopLogin;
};
