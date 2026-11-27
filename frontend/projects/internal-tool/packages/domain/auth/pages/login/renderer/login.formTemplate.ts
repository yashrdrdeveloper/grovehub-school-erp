import { DesktopLoginBrand } from '../templates/desktop/DesktopLogin';
import { MobileLoginBrand } from '../templates/mobile/MobileLogin';

export function getLoginBrandTemplate(mobile: boolean) {
  if (mobile) return MobileLoginBrand;
  return DesktopLoginBrand;
}
