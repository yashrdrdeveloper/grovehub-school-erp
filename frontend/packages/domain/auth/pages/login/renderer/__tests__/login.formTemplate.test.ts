import { DesktopLogin } from '../../templates/desktopLogin';
import { MobileLogin } from '../../templates/mobileLogin';
import { getLoginFormTemplate } from '../login.formTemplate';

describe('getLoginFormTemplate', () => {
  it('selects only the mobile form on mobile', () => {
    expect(getLoginFormTemplate(true)).toBe(MobileLogin);
  });

  it('selects only the desktop form otherwise', () => {
    expect(getLoginFormTemplate(false)).toBe(DesktopLogin);
  });
});
