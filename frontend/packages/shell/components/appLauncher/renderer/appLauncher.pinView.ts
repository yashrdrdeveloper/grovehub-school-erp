import { Pin, PinOff } from '@Components/icons';
import { cn } from '@Base/utilities/cn';
import styles from '../AppLauncherApp.module.css';

/** Pinning view state is derived, never mutated inside JSX. */
export const getAppLauncherPinView = (isPinned: boolean) => {
  if (isPinned) {
    return {
      pinLabel: __('Unpin'),
      PinIcon: PinOff,
      pinClass: cn(styles['catalog-app__pin'], styles['catalog-app__pin--active']),
    };
  }
  return {
    pinLabel: __('Pin'),
    PinIcon: Pin,
    pinClass: styles['catalog-app__pin'],
  };
};
