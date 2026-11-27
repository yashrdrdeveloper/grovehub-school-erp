import { Pin, PinOff } from '@Components/icons';
import { getAppLauncherPinView } from '../appLauncher.pinView';
describe('getAppLauncherPinView', () => {
  it('resolves the pin action with its accessible label', () => {
    expect(getAppLauncherPinView(true).PinIcon).toBe(PinOff);
    expect(getAppLauncherPinView(true).pinLabel).toBe('Unpin');
    expect(getAppLauncherPinView(false).PinIcon).toBe(Pin);
    expect(getAppLauncherPinView(false).pinLabel).toBe('Pin');
  });
});
