import { DesktopRbac } from '../templates/desktopRbac';
import { MobileRbac } from '../templates/mobileRbac';

export const getRbacTemplate = (isMobile: boolean) => {
  if (isMobile) return MobileRbac;
  return DesktopRbac;
};
