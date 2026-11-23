import { Pin, PinOff } from '@Components/icons';

/** Pinning view state is derived, never mutated inside JSX. */
export const getAppLauncherPinView = (isPinned: boolean) => {
  if (isPinned) {
    return {
      pinLabel: __('Unpin'),
      PinIcon: PinOff,
      pinClass: 'catalog-app__pin catalog-app__pin--active',
    };
  }
  return {
    pinLabel: __('Pin'),
    PinIcon: Pin,
    pinClass: 'catalog-app__pin',
  };
};
