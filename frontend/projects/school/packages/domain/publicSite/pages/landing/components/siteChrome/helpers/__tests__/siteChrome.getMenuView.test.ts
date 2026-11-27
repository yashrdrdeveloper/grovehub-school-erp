import { Menu, X } from '@Components/icons';
import { getSiteChromeMenuView } from '../siteChrome.getMenuView';
describe('getSiteChromeMenuView', () => {
  it('selects the matching icon and readable action', () => {
    expect(getSiteChromeMenuView(false).MenuIcon).toBe(Menu);
    expect(getSiteChromeMenuView(true).MenuIcon).toBe(X);
    expect(getSiteChromeMenuView(true).menuLabel).toBe('Close menu');
  });
});
