import { getAppLauncherPinAccessibleCopy } from '../appLauncher.getPinAccessibleCopy';

describe('getAppLauncherPinAccessibleCopy', () => {
  it('interpolates the application name in pin and unpin messages', () => {
    expect(getAppLauncherPinAccessibleCopy(false, 'Dashboard')).toEqual({
      accessibleName: 'Pin Dashboard',
      tooltip: 'Pin Dashboard to sidebar',
    });
    expect(getAppLauncherPinAccessibleCopy(true, 'Dashboard')).toEqual({
      accessibleName: 'Unpin Dashboard',
      tooltip: 'Unpin Dashboard from sidebar',
    });
  });
});
