import { DesktopAudit } from '../templates/desktopAudit';
import { MobileAudit } from '../templates/mobileAudit';

export const getAuditTemplate = (isMobile: boolean) => {
  if (isMobile) return MobileAudit;
  return DesktopAudit;
};
